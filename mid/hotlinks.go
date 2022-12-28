package mid

import (
	"context"
	"net/http"

	"gd9/prj3/kit/web"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// Hotlinks forbids requests from all Referer, Host, and Origin not on hosts whitelist.
func Hotlinks(hosts []string) web.Middleware {
	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.hotlinks")
			defer span.End()

			var err error

			if !matchHost(hosts, r.Host) {
				err = errors.New("hotlink : host mismatch : " + r.Host)
			}
			if !matchOrigin(hosts, r.Referer()) {
				err = errors.New("hotlink : referer mismatch")
			}
			o := r.Header.Get("Origin")
			if o != "" {
				if !matchOrigin(hosts, o) {
					err = errors.New("hotlink : origin mismatch")
				}
			}

			if err != nil {
				return web.NewRequestError(err, http.StatusForbidden)
			}
			return after(ctx, w, r)
		}
		return h
	}
	return m
}
