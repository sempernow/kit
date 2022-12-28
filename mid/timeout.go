package mid

import (
	"context"
	"net/http"
	"time"

	"gd9/prj3/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// ***********************************************************
// FAIL @ this Timeout middleware  : See LOG.md # 2020-07-29
// ***********************************************************

// Timeout sets the `context.WithTimeout()` to cancel next `web.Handler` after `t` milliseconds.
func Timeout(t time.Duration) web.Middleware {
	m := func(next web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			//ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.timeout")
			//defer span.End()
			ctx, cancel := context.WithTimeout(ctx, (t * time.Millisecond))
			defer cancel()

			ch := make(chan error, 1)
			go func() { ch <- next(ctx, w, r) }()

			select {
			case err := <-ch: // @ Response COMPLETED.
				return err
			case <-ctx.Done(): // @ Response CANCELLED.
				cancel() // Must here too, else fails to cancel many (most).
				<-ch     // Wait for graceful cancellation,
				// else RADICALLY increases failed (503) responses.
				return web.NewRequestError(ctx.Err(), http.StatusServiceUnavailable) // 503
			}
		}
		return h
	}
	return m
}

// TimeoutHandler adapts `http.TimeoutHandler` to run next `web.Handler` with timeout of `t` milliseconds.
// FAILs: Adapting to/from `http` pkg equivalents of `web.Handler`, but not of `web.ServeHTTP`.
func TimeoutHandler(t time.Duration) web.Middleware {

	// This is the function executed by the `wrapMiddleware` initialization function.
	m := func(next web.Handler) web.Handler {

		// This is the actual middleware function; a `web.Handler` wrapping the next in the chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.timeouthandler")
			defer span.End()

			_h := http.TimeoutHandler(WH2H(ctx, next), (t * time.Millisecond), "Timeout!.")
			//_h.ServeHTTP(w, r)
			return H2WH(_h)(ctx, w, r)
		}
		return h
	}
	return m
}

// WH2H ... `web.Handler` to `http.Handler`
func WH2H(ctx context.Context, wh web.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wh(ctx, w, r)
		},
	)
}

// H2WH ... `http.Handler` to `web.Handler`
func H2WH(h http.Handler) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.h2wh")
		defer span.End()

		h.ServeHTTP(w, r)
		return nil
	}
}
