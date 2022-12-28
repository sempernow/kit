package mid

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sempernow/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// **************************************************
// TODO: Migrate to https://github.com/uber-go/zap
// **************************************************

// Logger writes some request info to logs if request path
// does not contain any blacklisted string of []excludedPaths list.
//
//	Format: TraceID : (200) GET /foo/bar -> IP:Port (Latency)
func Logger(log *log.Logger, excludePaths ...string) web.Middleware {
	m := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.logger")
			defer span.End()

			// If the context is missing this value,
			// request the service to be shutdown gracefully.
			v, ok := ctx.Value(web.Key1).(*web.Values)
			if !ok {
				return web.NewShutdownError("context : missing web values")
			}

			err := before(ctx, w, r)

			// Filter what's logged ...
			var no bool
			for _, path := range excludePaths {
				//if r.URL.Path == path {
				if strings.Contains(r.URL.Path, path) {
					no = true
				}
			}
			if !no {
				// log.Printf("%s : (%d) : %s %s -> %s (%s)",
				// 	v.TraceID[:7]+"\u2026", // U+2026 HORIZONTAL ELLIPSIS
				// 	v.StatusCode,
				// 	r.Method,
				// 	r.URL.Path,
				// 	r.RemoteAddr,
				// 	time.Since(v.Now),
				// )
				log.Printf("(%d) : %s %s -> %s (%s)",
					v.StatusCode,
					r.Method,
					r.URL.Path,
					r.RemoteAddr,
					time.Since(v.Now),
				)
			}

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
