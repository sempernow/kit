// Package web extends the standard net/http pkg to provide a services framework.
package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/sempernow/kit/convert"

	"github.com/dimfeld/httptreemux/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key1 is used for protected set/get of Claims VALUEs from a `context.Context`.
const Key1 ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// RespTimeMax is app-wide max response time in milliseconds,
// measured from time of request arriving at its (first) endpoint handler;
// first in the middlewares chain.
const RespTimeMax = 9000

// Handler defines the per-request endpoint-handler type for this app framework.
// ... which BREAKS (???) the standard lib's `Handler` INTERFACE:
//
//	type Handler interface {
//	    ServeHTTP(ResponseWriter, *Request)
//	}
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// registered tracks handlers/routes added to our handle method,
// which invokes the http default servermux;
// a singleton used by the standard library for metrics and profiling.
// If a route is registered more than once, it would cause a panic.
var registered = make(map[string]bool)

// App is the entrypoint into our app and what configures our context object for http handlers.
// Add any configuration data/logic here.
type App struct {
	mux *httptreemux.ContextMux // https://github.com/dimfeld/httptreemux
	//... Handles trailing slashes per HTTP 301 Redirect.
	otmux    http.Handler
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an `App` value to handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	// The OpenTelemetry (OT) HTTP Handler (otmux) wraps this application's router (mux).
	// The OT handler starts initial span and annotates it with request/response info.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/

	mux := httptreemux.NewContextMux()

	return &App{
		mux:      mux,
		otmux:    otelhttp.NewHandler(mux, "req"),
		shutdown: shutdown,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP implements the `http.Handler` interface. It's the entry point for
// all http traffic and allows the opentelemetry mux to run first to handle
// tracing. The opentelemetry mux then calls the application mux to handle
// application traffic.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

// HandleDebug sets a handler function for a given HTTP method-path pair
// to the default http package server mux. /debug is added to the path.
func (a *App) HandleDebug(method string, path string, handler Handler, mw ...Middleware) {
	a.handle(true, method, path, handler, mw...)
} //... Does NOT respond the same as Handle; HTTP 404 and Binary output @ /debug/health

// Handle sets a handler function for a given HTTP method-path pair to the application server mux.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	a.handle(false, method, path, handler, mw...)
}

// handle applies the per-endpoint handler boilerplate and framework code.
func (a *App) handle(debug bool, method string, path string, handler Handler, mw ...Middleware) {
	if debug {
		// Track all the handlers that are being registered so we don't have
		// the same handlers registered twice to this singleton.
		if _, exists := registered[method+path]; exists {
			return
		}
		registered[method+path] = true
	}

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Start or expand a distributed trace.
		ctx := r.Context()
		ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, r.URL.Path)
		defer span.End()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, Key1, &v)

		// Sans timeout
		// // Call the wrapped handler functions.
		// if err := handler(ctx, w, r); err != nil {
		// 	a.SignalShutdown()
		// 	return
		// }

		// Add app-wide timeout
		ctx, cancel := context.WithTimeout(ctx, (RespTimeMax * time.Millisecond))
		defer cancel()

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {

			// OLD
			a.SignalShutdown()
			return

			// NEW
			// if validateShutdown(err) {
			// 	a.SignalShutdown()
			// 	return
			// }
		}

		// CANNOT HANDLE timeout HERE
		// Here, timout of two sequential requests generates (intractable) panic.
		// As bad, the downstream processes are NOT cancelled by this (below).
		// MUST HANDLE PER ENDPOINT (per process); note, e.g., sqlx takes ctx as arg.

		// ch := make(chan error, 1)
		// go func() { ch <- handler(ctx, w, r) }()

		// select {
		// case err := <-ch: // @ Response completion
		// 	if err != nil {
		// 		a.SignalShutdown()
		// 		return
		// 	}
		// case <-ctx.Done(): // @ Response timeout
		// 	cancel() // Must here too, else fails to cancel many (most).
		// 	<-ch     // Wait for graceful cancellation,
		// 	Respond(ctx, w, struct{ Err string }{
		// 		Err: fmt.Sprintf("%v", ctx.Err())},
		// 		http.StatusInternalServerError,
		// 	)
		// 	// The response "cancels" in that it returns upon context timeout,
		// 	// but the LOGger ... DOES NOT CANCEL
		// 	// ... logs only after the handler's (too lengthy) process ends.
		// 	// PWA : 2021/06/27 14:55:33.492667 errors.go:33: 0000000… : ERR : http: wrote more than the declared Content-Length
		// 	// PWA : 2021/06/27 14:55:33.492947 logger.go:43: 0000000… : (500) : HEAD /ops/sleep/10000 -> 127.0.0.1:2068 (10.0009579s)
		// 	return
		// }
	}

	// Add this handler for the specified verb and route.
	if debug {
		f := func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == method:
				h(w, r)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		http.DefaultServeMux.HandleFunc("/debug"+path, f)
		return
	}
	a.mux.Handle(method, path, h)
}

// Added ...

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
// See latest @ https://github.com/ardanlabs/service/blob/master/foundation/web/web.go
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}

// Redirect performs as http.Redirect(..), while fitting the context-based signature of this library.
func Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, url string, code int, msg ...string) error {
	if (code > 399) || (code < 300) {
		return errors.New("wrong http-status code for a redirect:" + convert.ToString(code))
	}
	bb := []byte{}
	if len(msg) > 0 {
		bb = []byte(strings.Join(msg, " : "))
	}
	w.Header().Set("Location", url)
	w.WriteHeader(code)
	return Respond(ctx, w, bb, code)
}
