package mid

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/sempernow/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// ************************************************************
// Handle gzip at response writer, not at request middleware.
// ************************************************************

// Extracted from: "Idiomatic golang net/http gzip transparent compression"
//  https://gist.github.com/CJEnright/bc2d8b8dc0c1389a9feeddb110f822d7
//  Author: CJEnright, 2020 version (MIT License)

// Gzip contains a chainable handler that efficiently (sync.Pool)
// performs gzip compression on the response body, adding the relevant header,
// but does so only if the request explicitly accepts gzip.
func Gzip() web.Middleware {
	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.gzip")
			defer span.End()

			// Abort lest request accepts gzip explicitly.
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				return after(ctx, w, r)
			}

			// MANY ISSUES:
			// 1. Crashes server @ HTTP 304
			// 2. Response mutated to HTTP 200 if HTTP 400
			// ...
			// ================================================================
			//  Gzip is CRASHING the SERVER on HTTP 304.
			//  Gzip writer wraps the zero-length body to some non-zero.
			//  Such is forbidden @ certain request/response.
			//  Also, middleware PRECEEDS primary handlers,
			//  so cannot be conditional on response body or status.
			// ================================================================

			w.Header().Set("Content-Encoding", "gzip")

			gz := gzPool.Get().(*gzip.Writer)
			defer gzPool.Put(gz)

			gz.Reset(w)
			defer gz.Close()

			return after(ctx, &gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
		}

		return h
	}

	return m
}

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
