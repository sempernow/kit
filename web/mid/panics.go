package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/sempernow/kit/web"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *log.Logger) web.Middleware {
	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.panics")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.Key1).(*web.Values)
			if !ok {
				return web.NewShutdownError("context : missing web values")
			}

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s :\n%s", v.TraceID, debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(ctx, w, r)
		}

		return h
	}

	return m
}
