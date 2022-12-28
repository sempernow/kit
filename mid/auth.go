package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/sempernow/kit/auth"
	"github.com/sempernow/kit/web"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

// TODO : Naming : alias the unhappy coincidence of the two words (auth/auth):
// Token : AUTHENTICATE against Token (creds)
// Roles : AUTHORIZE against roles

// ErrForbidden is returned when role of authenticated user is insufficient for action.
var ErrForbidden = web.NewRequestError(
	errors.New("claimant not authorized for that action : lacks required role"),
	http.StatusForbidden,
)

// ValidToken validates the Access token per request header `Authorization: Bearer <TOKEN>`.
// If invalid then responds with HTTP 403 if Refresh token-reference cookie is valid,
// else responds with HTTP 401.
func ValidToken(a *auth.Auth, rRefKey string) web.Middleware {
	// *********************************************************************************
	//  This scheme requires both JWT and cookies. Cookies are for Refresh-Token flag
	//  and one of several selectable CSRF mitigation methods.
	//
	// TODO: Pass in h.meta.Service (for Host, Origin). See HotlinksCSRF(..).
	// *********************************************************************************

	// This is the actual middleware function to be executed.
	m := func(after web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Start a span to measure just the time spent authenticating.

			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.authenticate")
			defer span.End()

			// Retrieve bearer token by parsing request header (`Authorization: Bearer <TOKEN>`).
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer TOKEN")
				return web.NewRequestError(err, http.StatusUnauthorized)
			} //... HTTP 401 if no such header.

			// Validate token (type unknown)
			// Cookie test is for Refresh token scheme; only upon expired Access token.
			// The API is cookie-less if client utilizes only the shorter-lived Access token
			// and if API does not use cookie-based (DoubleSubmitCookieMethod) CSRF mitigation.
			claims, errTkn := a.ValidateToken(parts[1])

			if errTkn != nil {
				// If token is invalid, then chk for Refresh token reference cookie;
				// cookie expiry matches that of its referenced token.
				_, err := r.Cookie(rRefKey)
				if err != nil {
					if err == http.ErrNoCookie { // If Refresh token reference not exist ...
						err = errors.Wrap(err, "token and refresh-reference cookie invalid")
						return web.NewRequestError(err, http.StatusUnauthorized)
					}
					return errors.Wrap(err, "reading refresh-token-reference cookie yet not expired")
				}
				// IIF valid Refresh token reference AND invalid Access token;
				// bearer has expired authorization yet is still authenticated,
				// which has no HTTP status code. Let's inform and remain semantically aligned ...
				err = errors.Wrap(errTkn, "valid refresh token but invalid access token")
				//w.Header().Set("Vary", "refresh-per-token")                 //... inform bearer of how.
				return web.NewRequestError(err, http.StatusForbidden) //... HTTP 403
				//... which FLAGs CLIENT to request new Access token by bearing valid Refresh token.
				// Note that clients are able to self test their token(s) for expiry directly, sans req/resp.
			}

			// Validate token type
			if claims.TokenType != auth.Access {
				err := errors.New("access token has invalid token-type claim")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// ADD Access token claims TO CONTEXT for downstream (per request) retrieval.
			ctx = context.WithValue(ctx, auth.Key1, claims)

			return after(ctx, w, r)
		}

		return h
	}

	return m
}

// ValidRoles validates that an AUTHENTICATED USER has Role-Based ACcess (RBAC);
// at least 1 role from a list of such; `auth.Claims{Roles: []string{auth.RoleUsrMbr}}`.
func ValidRoles(roles ...string) web.Middleware {
	m := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "kit.mid.authorize")
			defer span.End()

			claims, ok := ctx.Value(auth.Key1).(auth.Claims)
			if !ok {
				return errors.New("context : missing claims: Authorize called without/before Authenticate")
			}

			if !claims.Has(roles...) {
				return ErrForbidden
			}

			return after(ctx, w, r)
		}

		return h
	}

	return m
}
