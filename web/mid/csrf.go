package mid

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sempernow/kit/auth"
	"github.com/sempernow/kit/types/convert"
	"github.com/sempernow/kit/web"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// Modes of CSRF Mitigation
const (
	SansMitigation = iota
	CustomAJAXHeader
	SourceTargetHeaders
	DoubleSubmitCookie
	DomainLockedDouble
	HMACCookie
)

// CSRF mitigates Cross-Site Request Forgery attacks per mode.
// Each mode is a stateless OWASP-advised method for doing so.
// Multiple modes may be invoked per handler; modes are mutually orthogonal.
// See mitigate(..) for per-mode details.
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
//
//	USAGE: mid.CSRF(mid.DomainLockedDouble, "__Host-c", []string{"http://foo.com", "https://api.bar.xyz"}...)
func CSRF(mode int, cookieKey string, origins ...string) web.Middleware {

	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.csrf")
			defer span.End()

			err := mitigate(r, mode, cookieKey, origins...)
			if (mode == DoubleSubmitCookie) || (mode == DomainLockedDouble) {
				DeleteCSRFCookie(w, cookieKey)
			}
			if err != nil {
				return err
			}
			return after(ctx, w, r)
		}
		return h
	}
	return m
}

// DeleteCSRFCookie deletes that created (by client)
// in modes DoubleSubmitCookie or DomainLockedDouble.
func DeleteCSRFCookie(w http.ResponseWriter, key string) {
	csrf := &http.Cookie{
		Name:     key,
		Value:    "lolz",
		SameSite: auth.SameSiteMode,
		Expires:  time.Unix(0, 0), // 1970-01-01 00:00:00 +0000 UTC
		MaxAge:   -1,              // seconds
		Secure:   true,            // HTTPS only
		HttpOnly: false,           // allow JS access
		Path:     "/",
	}
	http.SetCookie(w, csrf)
}

// mitigate(..) per mode. The origins param applies only to
// SourceTargetHeaders (1) mode.
func mitigate(r *http.Request, mode int, key string, origins ...string) error {
	switch mode {
	case SansMitigation:
	case CustomAJAXHeader:
		// This solution is weak; relies on Fetch default settings/behavior.
		// Custom Request Header @ AJAX requests
		// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#use-of-custom-request-headers
		dbug(DBUG, "CSRF : CustomAJAXHeader: ", r.Header.Get("X-CSRF-Token"))
		if r.Header.Get("X-CSRF-Token") == "" {
			err := errors.New("missing x-csrf-token header")
			return web.NewRequestError(err, http.StatusForbidden)
		}
	case SourceTargetHeaders:
		// Verifying Source and Target Origins With Standard Headers
		// (Doubles as hotlink filter)
		// REF: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#verifying-origin-with-standard-headers
		//******************************************************************
		// Requests lacking Referer/Host/Origin headers will pass CSRF.
		// Such are typically legitimate; 3rd-party server-sent requests.
		// *****************************************************************
		var (
			err    error
			errStr string
		)
		// SOURCE (Origin||Referer)
		if !matchOrigin(origins, r.Referer()) {
			errStr = "host-referer"
		}
		if !matchOrigin(origins, r.Header.Get("Origin")) {
			if errStr != "" {
				errStr += "AND host-origin"
			}
		}
		// TARGET (Host) header must be sent in all HTTP/1.1 requests.
		// Golang unsets the Host request header; "promotes" it to `r.Host`.
		if !matchHost(origins, r.Host) {
			if errStr != "" {
				errStr += ", "
			}
			errStr = errStr + "host-target"
		}
		if errStr != "" {
			errStr = "csrf : mismatch : " + errStr
			err = errors.New(errStr)
		}

		dbug(DBUG, "CSRF : SourceTargetHeaders",
			" : Referer: ", r.Referer(),
			" : Host: ", r.Host,
		)
		if err != nil {
			return web.NewRequestError(err, http.StatusForbidden)
		}
	case DoubleSubmitCookie:
		// Double-submit Cookie Method of CSRF Attack Mitigation
		// Server-sent nonce cookie (value) also returned in the request (per Form or JSON payload).
		// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie
		fallthrough
	case DomainLockedDouble:
		// The DoubleSubmitCookie method using domain-locked cookie;
		// unreadable by cross-site attacker.

		// Extract CSRF token from its cookie
		c, err := r.Cookie(key)
		if err != nil {
			if err == http.ErrNoCookie {
				err := errors.Wrap(err, "csrf")
				return web.NewRequestError(err, http.StatusForbidden)
			}
			return errors.Wrap(err, "csrf token : reading cookie")
		}

		// Tee the response body for use here
		// while preserving it for downstream handler.
		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		r.Body = io.NopCloser(&buf)

		// Extract CSRF token from request body.
		type csrf = struct {
			CSRF string `json:"csrf"`
		}
		j := json.NewDecoder(body)
		var tkn csrf
		if err := j.Decode(&tkn); err != nil {
			return errors.Wrap(err, "csrf token : decoding body")
		}

		// Match the two extractions or forbid request.
		if tkn.CSRF != c.Value {
			err := errors.New("csrf token : body-cookie mismatch")
			return web.NewRequestError(err, http.StatusForbidden)
		}

		dbug(DBUG, "CSRF : DomainLockedDouble : token: ", tkn.CSRF)

	case HMACCookie:
		// HMAC-based Token Pattern : a DoubleSubmitCookie method
		// https://en.wikipedia.org/wiki/HMAC
		// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#hmac-based-token-pattern
		err := errors.New("csrf : this mitigation mode " + convert.ToString(HMACCookie) + " is not implemented")
		return web.NewRequestError(err, http.StatusForbidden)
	}

	return nil
}
