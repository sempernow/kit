package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sempernow/kit/convert"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// Content Types (MIME Types).
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
const (
	JSON        = "application/json"
	WEBMANIFEST = "application/manifest+json"

	GIF  = "image/gif"
	ICON = "image/x-icon"
	JPG  = "image/jpeg"
	PNG  = "image/png"
	SVG  = "image/svg+xml"
	//SVG  = "application/xml" //... FAIL @ browser.
	//... CloudFront resets content-type to this if file's 1st line is HTML comment.
	WEBP = "image/webp"

	CSS  = "text/css"
	HTML = "text/html"
	JS   = "text/javascript"

	WOFF = "font/woff"

	BOGUS       = "bogus"
	MALFORMED   = "malformed"
	UNSUPPORTED = "unsupported"
)

const cspSEP = "; "

// RespStatus contains code and text of an HTTP response.
type RespStatus struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Status returns RespStatus of an HTTP response code.
func Status(status int) RespStatus {
	return RespStatus{
		Code: status,
		Text: http.StatusText(status),
	}
}

// Resource contains all the parameters of a web-server resource.
type Resource struct {
	Key     string
	Content []byte
	Ctype   string
	Etag    string
	SRI     string
	Ext     string
	Mtime   time.Time
	Ctime   time.Time
	Model   bool
	Gz      bool
	Err     error
	Code    int
}

// CSP contains lists of Content Security Policy (CSP) sources.
type CSP struct {
	ConnectSrc, ScriptSrc, StyleSrc, FrameSrc, FontSrc, ImgSrc, ObjectSrc, DefaultSrc []string
}

// Respond is the nominal Response function; enforces strictest CSP.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	return NewResponse(CSP{})(ctx, w, data, statusCode)
}

// RespondError is the nominal RespondError function; enforces strictest CSP.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	return NewResponseError(CSP{})(ctx, w, err)
}

// RespondFoo is a Response function that loosens Content Security Policy (CSP)
// to allow connections, scripts and such, each from a whitelist of sources.
// Written as template for custom response function(s) per service or endpoint.
var RespondFoo = NewResponse(CSP{
	// CSP rejects all sources except origin server ('self') lest whitelisted hereby.
	ConnectSrc: []string{
		"https://foo.foo",
	},
	ScriptSrc: []string{
		"https://foo.foo",
	},
	StyleSrc: []string{
		"'unsafe-inline'",

		"https://foo.foo",
	},
	DefaultSrc: []string{
		"https://foo.foo",
	},
})

// Response is the per-request HTTP response function.
type Response func(context.Context, http.ResponseWriter, interface{}, int) error

// NewResponse closes over CSP sources (whitelists), returning a Response function.
// TODO : Add custom marshaller args to handle enums handling of whatever types
// https://stackoverflow.com/questions/38897529/pass-method-argument-to-function
// https://rotational.io/blog/marshaling-go-enums-to-and-from-json/
func NewResponse(cc CSP) Response {
	csp := struct {
		ConnectSrc,
		FrameSrc,
		FontSrc,
		ScriptSrc,
		StyleSrc,
		ImgSrc,
		ObjectSrc,
		DefaultSrc string
	}{
		ConnectSrc: strings.Join(cc.ConnectSrc[:], " "),
		FrameSrc:   strings.Join(cc.FrameSrc[:], " "),
		FontSrc:    strings.Join(cc.FontSrc[:], " "),
		ScriptSrc:  strings.Join(cc.ScriptSrc[:], " "),
		StyleSrc:   strings.Join(cc.StyleSrc[:], " "),
		ImgSrc:     strings.Join(cc.ImgSrc[:], " "),
		ObjectSrc:  strings.Join(cc.ObjectSrc[:], " "),
		DefaultSrc: strings.Join(cc.DefaultSrc[:], " "),
	}

	/* CSP Header:
	   Content-Security-Policy: connect-src 'self' https://github.com https://amazon.com https://paypal.com https://google.com https://api.authorize.net https://apitest.authorize.net https://js.authorize.net https://jstest.authorize.net; script-src 'self' https://js.authorize.net https://jstest.authorize.net; style-src 'self' 'unsafe-inline' https://js.authorize.net https://jstest.authorize.net; font-src 'self' data:; default-src 'self' https://api.authorize.net https://apitest.authorize.net
	*/

	// The HTTP response function.
	return func(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

		ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.web.response")
		defer span.End()

		// Set status code for request-logger middleware.
		// If context is missing Value, then request service shutdown gracefully.
		v, ok := ctx.Value(Key1).(*Values)
		if !ok {
			return NewShutdownError("context : missing web values")
		}
		v.StatusCode = statusCode

		var (
			bb      []byte
			ctype   string
			isGz    bool
			err     error
			nocache bool

			cspConnectSrc = csp.ConnectSrc
			cspFrameSrc   = csp.FrameSrc
			cspScriptSrc  = csp.ScriptSrc
			cspFontSrc    = csp.FontSrc
			cspObjectSrc  = csp.ObjectSrc
			cspImgSrc     = csp.ImgSrc
			cspStyleSrc   = csp.StyleSrc
			cspDefaultSrc = csp.DefaultSrc

			// If `Page`, then almost always INCLUDE BODY in response,
			// regardless of HTTP-status code; except @ HTTP 204 and 304 ...
			sansBody = (statusCode == http.StatusNoContent) ||
				(statusCode == http.StatusNotModified)
		)

		switch sansBody {
		case true:
			// --------------------------------------------------------------------
			// @ Set required header @ 304 (if would send @ 200)

			if statusCode == http.StatusNotModified {
				switch as := data.(type) {
				case *Resource:
					if !(as.Etag == "" || as.Etag == BOGUS) {
						w.Header().Set("Etag", (`W/"` + as.Etag + `"`))
					}
				}
			}
		case false:
			// --------------------------------------------------------------------
			// @ Convert resource to bytes (body), and othewise prep response.

			// ====================================================================
			//  TODO: Move these switch/case blocks to their respective callers;
			//        all callers send a `*Resource`; callers do response prep.
			// ====================================================================
			switch as := data.(type) {

			case *Resource:
				ctype = as.Ctype
				bb = as.Content
				isGz = as.Gz
				if !(as.Etag == "" || as.Etag == BOGUS) {
					w.Header().Set("Etag", (`W/"` + as.Etag + `"`))
				} else {
					if !as.Mtime.IsZero() {
						w.Header().Set("Last-Modified",
							LastModified(as.Mtime),
						)
					} else {
						w.Header().Set("Last-Modified",
							LastModified(time.Now().UTC()),
						)
					}
				}
				// switch as.Ext {
				// case "json", "webmanifest":
				// 	cspDefaultSrc = csp.DefaultSrc + " " + as.SRI
				// case "js":
				// 	cspScriptSrc = csp.ScriptSrc + " " + as.SRI
				// case "css":
				// 	cspStyleSrc = csp.StyleSrc + " " + as.SRI
				// }

				if ctype == JSON {
					nocache = true
				}

			// ********************
			//   DEPRICATED CASEs
			// ********************

			case []byte:
				ctype = HTML
				bb = as

			case *bytes.Buffer: //... ??? ... coded @ prior version
				if !sansBody {
					w.Header().Set("Content-Type", HTML+"; charset=UTF-8")
					w.Header().Set("X-Content-Type-Options", "nosniff")
					if !strings.Contains(w.Header().Get("Content-Encoding"), "gzip") {
						w.Header().Set("Content-Encoding", "identity")
					}
				}
				w.WriteHeader(statusCode)
				_, err := as.WriteTo(w) // Stream
				//... ??? This scheme is questionable and out-of-band;
				// not in http pkg; does it close the response body?
				return err

			case io.Reader:
				ctype = HTML
				bb = convert.ReaderToBytes(as)

			case string:
				ctype = JSON
				bb = []byte(as)

			// ********************************
			//   API or RespondError(..) CASE
			// ********************************

			default: // Struct
				ctype = JSON
				bb, err = json.Marshal(as)
				if err != nil {
					return err
				}
				nocache = true
			}
		}

		// // Cache only the PWA-service responses; never that of AOA or API.
		// if ctype == JSON {
		// 	nocache = true
		// }

		// /*******************************
		//   TEST/DEBUG : Cache HTML only
		// *******************************/
		// nocache = true
		// if ctype == HTML {
		// 	nocache = false
		// }

		// --------------------------------------------------------------------
		// @ Conditionally gzip

		// TODO: integrate this into final writer; w.Write(bb)
		// (Middleware has no access to this post-processed `bb`)
		// if !isGz && len(bb) > 1024 {
		// 	isGz = true
		// 	bb, err = gz.Write(bb)
		// 	if err != nil {
		// 		return err
		// 	}
		// }//... REQUIRES MODS @ app-layer TESTING, which examine response body.

		// --------------------------------------------------------------------
		// @ Add certain headers only if non-zero body

		// @ References:
		// 	Headers  https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers
		// 	Status   https://httpstatuses.com/

		switch len(bb) > 0 {
		case true:
			// ----------------------------------------------------------------
			// @ CSP / SRI

			/****************************************************************************************
			  CSP References:
			  CanIUse  https://caniuse.com/#search=Content-Security-Policy
			  MDN  https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
			*****************************************************************************************/
			//... 'unsafe-inline' setttings are required by AccessUI.js scheme
			// of Authorize.net for PCI-DSS compliance, ironically.

			//if ctype == HTML || ctype == JS || ctype == CSS {
			//if ctype == JS || ctype == CSS || ctype == JSON {
			//if ctype == JS || ctype == CSS {
			if ctype == HTML {
				csp := []struct {
					self, more string
				}{
					{
						self: "connect-src 'self' ",
						more: cspConnectSrc,
					}, {
						self: "frame-src 'self' ",
						more: cspFrameSrc,
					}, {
						self: "script-src 'self' ",
						more: cspScriptSrc,
					}, {
						self: "style-src 'self' ",
						more: cspStyleSrc,
					}, {
						self: "font-src 'self' ",
						more: cspFontSrc,
					}, {
						self: "object-src 'self' ",
						more: cspObjectSrc,
					}, {
						self: "img-src 'self' ",
						more: cspImgSrc,
					},
				}
				var sources string
				for _, v := range csp {
					if v.more != "" {
						sources = sources + v.self + v.more + cspSEP
					}
				}
				w.Header().Set("Content-Security-Policy",
					sources+"default-src 'self' "+cspDefaultSrc,
				)
			}

			if !strings.Contains(ctype, "image/") {
				ctype += "; charset=UTF-8"
			}

			w.Header().Set("Content-Type", ctype)
			w.Header().Set("X-Content-Type-Options", "nosniff")

			switch isGz {
			case true:
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Content-Length", convert.IntToString(len(bb)))
			case false:
				w.Header().Set("Content-Encoding", "identity")
				w.Header().Set("Content-Length", convert.IntToString(len(bb)))
			}

		case false:
			w.Header().Set("Content-Length", "0")
		}

		if statusCode == 401 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Access the web app APIs"`)
			/****************************************************************************************
			  Types : Basic(*)|Bearer|Digest(*)|HOBA|Mutual|Negotiate|OAuth|SCRAM=SHA-{1,256}|vapid
							* Types that trigger browsers' Login-widget popup (blocking).

			  REF: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate
				   https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml
			****************************************************************************************/
		}

		if statusCode > 499 {
			nocache = true
		}

		// ------------------------------------------------------------------------
		// @ Add certain response headers REGARDLESS

		/*************************************************************************
		  If HTTP 304, then MUST INCLUDE certain HEADERS (if would send @ 200):
		  Cache-Control, Content-Location, Date, ETag, Expires, and Vary.
		*************************************************************************/
		switch nocache {
		case true:
			/************************************************************************
			  Absent this header, some browsers cache the JSON format
			  of a Page model (requested per content negotiation; same URL),
			  thereby replacing the (older) HTML payload,
			  thereafter rendering JSON instead of HTML on page reload.
			************************************************************************/
			w.Header().Set("Cache-Control", "no-store, no-transform")
			/************************************************************************
			  ... But this response header prevents Nginx from caching.
			  We want server to cache, but client (browser per se) to not cache.
			  Solution is to always cache (server side) and then modify this header
			  post cache, at reverse proxy. Client-side cache is then purely
			  per ServiceWorker/Cache Web APIs,
			  thereby removing browser quirks categorically.

			  Nginx lacks logic required of our app; per content negotiation, etc.
			************************************************************************/
		case false:
			w.Header().Set("Cache-Control", "public, max-age=31536000, no-transform, immutable")
		}

		// ------------------------------------------------------------------------
		// @ Send it

		w.WriteHeader(statusCode)
		if _, err := w.Write(bb); err != nil {
			return err
		}
		return nil
	}
}

// ResponseError is the per-request HTTP response function for errors.
type ResponseError func(context.Context, http.ResponseWriter, error) error

// NewResponseError closes over CSP sources (whitelists), returning a ResponseError function.
func NewResponseError(csp CSP) ResponseError {

	response := NewResponse(csp)

	// The HTTP error response function.
	return func(ctx context.Context, w http.ResponseWriter, err error) error {
		// If error is of type `*Error`,
		// then handler returned specific HTTP status code and error.
		if webErr, ok := errors.Cause(err).(*Error); ok {
			er := ErrorResponse{
				Error:  webErr.Err.Error(),
				Fields: webErr.Fields,
			}
			if err := response(ctx, w, er, webErr.Status); err != nil {
				return err
			}
			return nil
		}

		// Else handler sent an arbitrary error value,
		er := ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		if err := response(ctx, w, er, http.StatusInternalServerError); err != nil {
			return err
		} //... so respond with HTTP 500.
		return nil
	}
}

// ==================
//  HELPERs
// ==================

// LastModified is a header-helper function that returns the properly formatted value
// to fit the HTTP "Last-Modified: <LastModified>" header: `Thu, 20 Aug 2020 18:26:03 GMT`.
//
//	Usage: w.Header().Set("Last-Modified", LastModified(time.Now().UTC()))
func LastModified(x time.Time) string {
	return fmt.Sprintf("%s, %02d %s %02d %02d:%02d:%02d GMT",
		x.Weekday().String()[:3], x.Day(), x.Month().String()[:3],
		x.Year(), x.Hour(), x.Minute(), x.Second(),
	)
}

// IfModifiedSince is a header-helper function that tests if `subject` is modified `since`.
// Returns `true` is `subject` is newer than `since`.
// Returns `true` on either parse error.
// Time-string format is that of such headers (RFC1123):
//
//	"Thu, 20 Aug 2020 18:26:03 GMT"
func IfModifiedSince(subject, since string) bool {
	su, err := time.Parse(time.RFC1123, subject)
	if err != nil {
		return true
	}
	si, err := time.Parse(time.RFC1123, since)
	if err != nil {
		return true
	}
	return su.After(si)
}
