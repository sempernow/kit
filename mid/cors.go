package mid

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/sempernow/kit/str"
	"github.com/sempernow/kit/web"

	"go.opentelemetry.io/otel/trace"
)

// Are we doing it wrong? https://www.keycdn.com/blog/cors-cdn

/**********************************************************************************
* Some browsers erroneously report CORS requests that fail for unrelated reasons,
* e.g., response timeout, as a "Cross-Origin Request Blocked ...".
* Firefox does that, whereas Chrome reports the failed preflight (OPTIONS).
**********************************************************************************/

// CORS handles Cross-Origin Resource Sharing requests;
// allowed origins and methods are settable per endpoint,
// and closed over per declaration.
//   - methods GET, HEAD, and OPTIONS are allowed regardless.
func CORS(origins, methods []string) web.Middleware {

	if len(origins) == 0 {
		origins = append(origins, "*")
	}
	methods = append(methods, []string{"GET", "HEAD", "OPTIONS"}...)
	allow := strings.Join(str.Unique(methods), ",")

	match := func(origins []string, r *http.Request) string {
		for _, o := range origins {
			if r.Header.Get("Origin") == o {
				return r.Header.Get("Origin")
			} else {
				if o == "*" {
					return r.Header.Get("Origin")
				}
			}
		}
		return ""
	}

	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.cors")
			defer span.End()

			// --------------------------------------------------------------------
			// @ Preflight request

			if r.Method == "OPTIONS" { //... may be EITHER same-site OR cross-origin request.
				//... https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS

				// if !validateOrigin(r, meta) {
				// 	err := errors.New("cors origin not moted")
				// 	return web.NewRequestError(err, http.StatusBadRequest) // 400 (Nginx 444)
				// }

				// if method := r.Header.Get("Access-Control-Request-Method"); method != "" {
				// 	// Browsers AUTOMATICALLY send preflight @ certain cross-origin requests.
				// 	// https://developer.mozilla.org/en-US/docs/Glossary/preflight_request

				// 	allow := []string{}
				// 	allow = append(allow, strings.Split(allowed, ",")...)
				// 	//allow = append(allow, []string{"PUT", "DELETE"}...)

				w.Header().Set("Access-Control-Allow-Methods", allow)
				//... wildcard ("*") not allowed @ fetch w/ credentials
				// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods"

				// TODO: Read Origin header, validate against approved list (meta.Service),
				// and then set origin to Origin if allowing.

				if mo := match(origins, r); mo != "" {
					w.Header().Set("Access-Control-Allow-Origin", mo)
					//w.Header().Set("Access-Control-Allow-Origin", "*")
					// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
					//... EITHER one origin, or all (*)
					//    Client declares per value of request header `Origin: ...`.
					//    https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
				} else {
					if len(origins) > 0 {
						w.Header().Set("Access-Control-Allow-Origin", origins[0])
					}
				}
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				// List of Allowable Request Headers:
				// @ "Unsafe" requests : https://javascript.info/fetch-crossorigin#safe-requests
				// If request per Fetch API with `credentials = "include"` option,
				// then wildcard (*) here is read by browsers as a literal (meaningless).
				w.Header().Set("Access-Control-Allow-Headers",
					//"*",
					//"Authorization",
					//... sufficient @ API per fetch params @ Auth module (auth.js)
					//"Authorization,Content-Type",
					//... sufficient @ AOA POST (Why req declaring?)
					"Authorization,Content-Type,Cache-Control,If-Modified-Since,X-CSRF-Token",
					//... sufficient @ API per Net module (net.js)
					//"Authorization,Cookie",
				)
				w.Header().Set("Access-Control-Max-Age", "86400")
				//... after which client must preflight again.
				//... HAS NO AFFECT (@ Firefox/Brave).

				w.Header().Set("X-CORS-Test",
					"PREFLIGHT @ "+time.Now().UTC().Format(time.RFC3339),
				)
				//... https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow
				w.WriteHeader(http.StatusOK)
				//... HTTP 200 @ CORS preflight response signals permissive.

				return after(ctx, w, r)
			}

			// --------------------------------------------------------------------
			// @ Main request

			ok := (r.Header.Get("Origin") == "")
			if ok { //... @ not CORS
				return after(ctx, w, r)
			} //... All browsers add Origin header to all CORS requests.

			if mo := match(origins, r); mo != "" {
				dbug(DBUG, "CORS : origins match :: ALLOWing")
				// To allow, send this header.
				w.Header().Set("Access-Control-Allow-Origin", mo)
				// To send Cookie, WWW-Authentication headers:
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				// ... browser sends only on fetch setting: "credentials: true"

			} else {
				if mo != "*" { //... if response varies per origin, then ...
					w.Header().Set("Vary", "Origin")
				} //... if request denied, give client a hint of how response varies.
			}

			// @ "Safe" requests : https://javascript.info/fetch-crossorigin#safe-requests
			// (@ "Unsafe" requests, this is handled @ preflight per Access-Control-Allow-Headers)
			//w.Header().Set("Access-Control-Expose-Headers", "*")
			w.Header().Set("Access-Control-Expose-Headers", "*")
			//... whitelist response headers exposed to client (JS).
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
			// Why: https://stackoverflow.com/questions/25673089/why-is-access-control-expose-headers-needed

			return after(ctx, w, r)
		}
		return h
	}
	return m
}

// ==================================
// r.Header @ Nginx as Reverse Proxy
// ==================================
// "Content-Length": ["0"]
// "Accept-Encoding": ["gzip, deflate"]
// "Referer": ["http://localhost/app/login"]
// "Authorization": ["Basic REABFyFXSlUJRlwETQddWzpib3o5NyN5"]
// "Cache-Control": ["no-cache"]
// "User-Agent": ["Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0"]
// "Accept": ["*/*"]
// "Accept-Language": ["en-US,en;q=0.5"]
// "Origin": ["http://localhost"]
// "Cookie": ["_na=ae931f95fac76f2bbf49ba2a2ab227e567f8e88c587e0f227b197a47b0a34d87b12bc72e43e5ac2dee43b82eda90f18365167f0e5562512a62dca06d422aed31; _nb=dccd5f8894057ccdb971f4b7b63d5b9709a58c7644b1551f8a4a01dd32a513b5460e46de8facd296cf4687a52222d19f7b326a54f5a942e70be13ac4eb182162; _nc=468f800d66256c63c5e1ba66d3053ac74f880bbd6dc610c9a82440eb41f352b72a17fd60977aa671de4ff0361b6e6ba43876065b2af7699efa637c36b37cb3ed"]
// "Pragma": ["no-cache"]
// "Connection": ["close"]
