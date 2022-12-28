// Package mid provides all app-layer middleware functions.
package mid

import (
	"context"
	"net/http"

	"github.com/sempernow/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// ExampleMiddleware(..)
// Each middleware function abides the same pattern and signature, but for args.
// Each returns a closure that returns a chainable, per-request handler.
// This transforms the otherwise nested syntax of middlewares into a CSV list,
// with execution abiding list order; all executing prior to the base handler.
// Such middleware may be appended per application service, or per route:
//
//	svc := web.NewApp(shutdown, mid.M1(bar), mid.M2(foo, whatever), ...)
//	svc.Handle(method, path, x.hndlrY, mid.M9(423), mid.M5(true), ...)
func ExampleMiddleware(foo int, bar string, whatever ...interface{}) web.Middleware {

	//... Close over foo, bar, and whatever here ...

	// Define the chainable middleware function returned on init.
	m := func(next web.Handler) web.Handler {

		// Define the per-request web.Handler
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.example")
			defer span.End()

			// ... Doings of this middleware ...

			return next(ctx, w, r)
		}
		return h
	}
	return m
}
