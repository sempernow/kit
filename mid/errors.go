package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/sempernow/kit/web"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel/trace"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *log.Logger) web.Middleware {

	logErr := func(v *web.Values, trace bool, e interface{}) {
		if trace {
			// 0000000… : ERR : ...
			log.Printf("%s : ERR : %v", v.TraceID[:7]+"\u2026", e)
		} else {
			log.Printf("ERR : %v", e)
		}
	}

	m := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.errors")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.Key1).(*web.Values)
			if !ok {
				return web.NewShutdownError("context : missing web values")
			}

			// Run the rest of the handler chain, catching any propagated errors.
			if err := before(ctx, w, r); err != nil {

				// Log the error w/out trace.
				trc := false
				if true {
					if webErr, ok := errors.Cause(err).(*web.Error); ok {
						// web.ErrorResponse is an abomination.
						// Subkey everything under the one Error key
						// Here, for the logger AND in the RespondError(..).
						erx := web.ErrorResponse{
							Error:  webErr.Err.Error(),
							Fields: webErr.Fields,
						}
						logErr(v, trc, erx)
					} else {
						logErr(v, trc, err)
					}
				} else {
					logErr(v, trc, err)
				}

				// Respond to the error.
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// The error has been handled so we can stop propagating it.
			return nil
		}

		return h
	}

	return m
}
