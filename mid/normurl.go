package mid

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"gd9/prj3/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// WHY ?

// NormURL rewrites the request URL by stripping prefix
// WIP
func NormURL(params string) web.Middleware {
	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.normalurl")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			_, ok := ctx.Value(web.Key1).(*web.Values)
			if !ok {
				return web.NewShutdownError("context : missing web values")
			}

			params := web.Params(r)
			prefix := params["type"] + "/" + params["path"]

			//if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			if params["path"] == "pub" {

				p := strings.TrimPrefix(r.URL.Path, prefix)

				r2 := new(http.Request)

				*r2 = *r

				r2.URL = new(url.URL)

				*r2.URL = *r.URL

				r2.URL.Path = p

				return after(ctx, w, r2)
			}
			return after(ctx, w, r)
		}

		return h
	}

	return m
}
