// Package auth provides authentication and authorization support.
package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid" // FORKED from github.com/satori

	//"github.com/dgrijalva/jwt-go" // v3.x; NO LONGER MAINTAINED
	"github.com/golang-jwt/jwt" // v3.x; Community-maintained fork.

	// TODO : JWT @ https://github.com/lestrrat-go/jwx
	//	"pkg.go.dev/github.com/lestrrat-go/jwx"
	//	"github.com/lestrrat-go/jwx"

	"github.com/pkg/errors"
)

// RBAC per Claims.Roles
const (
	// RBAC : App-wide

	RoleAppOps = "AOPS" // ADMIN
	RoleAppMod = "AMOD"
	RoleAppBot = "ABOT" // @ all mv_users_cfg records

	// RBAC : Per service

	RoleSvcADM = "SADM"
	RoleSvcAOA = "SAOA"
	RoleSvcAPI = "SAPI"
	RoleSvcPWA = "SPWA"

	// RBAC : Clients

	RoleUsrVis = "UVIS"
	RoleUsrMbr = "UMEM" // USER
	RoleUsrMod = "UMOD"
	//RoleUsrAPI  = "UAPI"

	RoleGrpMbr  = "GMEM"
	RoleGrpMod  = "GMOD"
	RoleGrpLead = "GLEA"

	RoleChnHost = "HOST"
	RolePtnrMem = "PMEM"
	RolePtnrApp = "PAPP"

	// ------------------
	// PRIOR : DEPRICATED
	//RoleUser  = "USER"
	//RoleAdmin = "ADMIN"
)

// Auth Modes
const (
	SignUp      = "SU"
	BasicAuth   = "BA"
	ObfuscateBA = "OB"
	DigestAuth  = "DA"
	OAuth2      = "OA"
	WebAuthn    = "WA"
)

// Access Mode values identify the issuer/claimant mode of authorization;
// the authentication endpoint/method (reference) by which the token was granted.
const (
	// Access Mode : Web-app user (tokens and cookies; pair of pairs)
	PWA = "fc160f10-5683-4cdf-b0b4-8ddee3a631bc"
	// Access Mode : API user (access token per api key)
	API = "5066494a-22ef-47ee-903b-43782244c08c"
	// Access Mode : AOA (claim existing user account per api key)
	AOA = "deb726dd-06a8-466c-b8b0-72082051c210"
)

const (
	// IssuedAt Offset [seconds]; allow for server-clock discrepencies.
	IssuedAtOffset = 10 // Else susceptible to ERR: "Token used before issued"

// ****************************************************************************
// ERR: "Token used before issued"
//
// This is a known issue of dgrijalva/jwt-go pkg;
// a fail mode @ token validation (`iat` claim)
// when server/clock thereof differs from that of issuer.
//
// FIX: Use IssuedAt(..) and IssuedAtNow() to set the `iat` claim.
//
//	https://github.com/dgrijalva/jwt-go/issues/383
//	https://github.com/golang-jwt/jwt/issues/98
//
// The newer, maintained fork (golang-jwt/jwt) was adopted after this fix,
// but a quick look at relevant (seemingly identical) code suggests same bug.
//
// TODO : Upgrade to lestrrat-go/jwx pkg.
// ****************************************************************************
)

// Token Type
const (
	/**************************************************************************
	Refresh token authorizes only a refresh request upon expired access token.
	This is the long-lived bearer token; it and its paired reference cookie,
	together, proxy user authentication. To validate,
	the iss claim of token must match the hash of its reference-cookie value.

	UPDATE: analogous scheme @ https://github.com/lestrrat-go/cert_bound_sts_server (2022-07-31)

	Cookie value:    ISS_REF_VAL         ... is not visible to attacker.
	Token ISS claim: SHA256(ISS_REF_VAL) ... is visible to attacker.
	**************************************************************************/
	Refresh = "R"

	/**************************************************************************
	Access token authorizes access to application resources.
	This is the short-lived bearer token; it alone proxies authorization.
	Its paired reference cookie is but a proxy for its own expiry upon refresh.
	**************************************************************************/
	Access = "A"
)

// Token TTLs utilized at data layer and app-wide testing.
// Services (app layer) declare these params per config on init;
// they are resettable (per environment vars) sans rebuild.
//
//	@Seconds: AccessTTL.Seconds(), @Epoch: time.Now().Add(AccessTTL).Unix()
const (
	RefreshTTL = 180 * 24 * time.Hour
	//RefreshTTL = 7 * time.Minute

	AccessTTL = 1 * time.Hour
	//AccessTTL = 2 * time.Minute

) //... auth-endpoint dynamics expect token-cookie pairs to have matching TTLs.

// Auth keys; client/server paramaters.
const (

	// Key names of client-side auth store
	KeyTknAccess  = "a"
	KeyTknRefresh = "r"
	//... used server side only at (obsolete) unit tests

	// Cookie param(s)

	SameSiteMode = http.SameSiteStrictMode

	// Token-reference cookies

	// Set-Cookie: __Host-a=VAL; path=/; SameSite: strict; Secure; HttpOnly;
	KeyRefAccess = "__Host-a"
	// Set-Cookie: __Host-r=VAL; path=/; SameSite: strict; Secure; HttpOnly;
	KeyRefRefresh = "__Host-r"

	// OAuth and CSRF-token cookie keys

	// Set-Cookie: __Host-o=VAL; path=/; SameSite: strict; Secure;
	KeyOA = "__Host-o" // FAIL @ Chrome/Brave @ localhost
	// Set-Cookie: __Host-c=VAL; path=/; SameSite: strict; Secure;
	KeyCSRF = "__Host-c"
)

// TokenPair contains the auth-token response body;
// Access and Refresh tokens, authentication mode,
// and the provider if per OAuth mode.
type TokenPair struct {
	A        string `json:"a,omitempty"`
	R        string `json:"r,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key1 is used for protected store/retrieve of Claims VALUEs from a `context.Context`.
const Key1 ctxKey = 1

// Claims represents the authorization claims of a bearer token (JSON Web Token; JWT).
// https://tools.ietf.org/html/rfc7519#section-4.1 | RFC7519
type Claims struct {
	jwt.StandardClaims
	Roles     []string `json:"roles"`
	TokenType string   `json:"tokenType"`
	Key       string   `json:"key"`
	//RemoteOrigin string   `json:"remote_origin,omitempty"`
	//ChnID        string   `json:"chn_id,omitempty"`
	//ChnSlug      string   `json:"chn_slug,omitempty"`
	//Internal     interface{} `json:-`
	//... Internal use only; excluded from the JWT itself.
}

// Issuer identifies authentication endpoint (e.g., URI) whereof token is issued;
//
//	Standard Claims : Issuer (iss) value.
func Issuer(principal string) string {
	return uuid.NewV5(uuid.NamespaceURL, principal).String()
}

// IssuedAt injects an offset to compenstate for server-clock discrpencies.
// Absent this, token validation may fail ("Token used before issued")
// if issuer and test run on different nodes.
// This occurred at private node that failed to synch with network time
// due to Ubuntu's newer time-synch utility; npt v. timedatectl.
//
//	Standard Claims : IssuedAt (iat) value.
func IssuedAt(now *time.Time) int64 {
	t := now.Add((-IssuedAtOffset) * time.Second)
	return t.Unix()
}
func IssuedAtNow() int64 {
	return time.Now().Add((-IssuedAtOffset) * time.Second).Unix()
}

// Has returns true if the claims has at least one of the provided roles.
func (c Claims) Has(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

// PubKeyLookup defines the signature of a function to lookup public keys.
// Asymmetric-signed JWTs contain only a reference (`kid`) to their public key.
// This is in their header segment: {"kid": <KID>, "typ": "JWT", "alg": "HS256"}.
//
// In a production system, a key id (KID) is used to retrieve the correct
// public key to parse a JWT for auth and claims. A key lookup function is
// provided to perform the task of retrieving a KID for a given public key.
//
// A key lookup function is required for creating an Authenticator.
//
// * Private keys should be rotated. During the transition period, tokens
// signed with the old and new keys can coexist by looking up the correct
// public key by KID.
//
// * KID to public key resolution is usually accomplished via a public JWKS
// endpoint. See https://auth0.com/docs/jwks for more details.
type PubKeyLookup func(kid string) (*rsa.PublicKey, error)

// JWKS (JSON Web Key Set) function returns a PubKeyLookup
func JWKS(activeKID string, publicKey *rsa.PublicKey) PubKeyLookup {
	f := func(kid string) (*rsa.PublicKey, error) {
		if activeKID != kid {
			return nil, fmt.Errorf("unrecognized key id %q", kid)
		}
		return publicKey, nil
	}

	return f
}

// Auth contains the `New(..)` authenticator's parameters; used to authenticate clients.
type Auth struct {
	privateKey          *rsa.PrivateKey
	activeKID           string
	algorithm           string
	keyFunc             PubKeyLookup
	parser              *jwt.Parser
	cookieKeyTknRefresh string
	cookieKeyTknAccess  string
}

// New creates an authenticator (`*Auth`) used to generate a token (JWT)
// for a set of user claims and recreate the claims by parsing the token.
// It will error if:
//   - The private key is nil
//   - The public key func is nil.
//   - The key ID is blank.
//   - The specified algorithm is unsupported.
func New(privateKey *rsa.PrivateKey, activeKID, algorithm string, lookup PubKeyLookup) (*Auth, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if activeKID == "" {
		return nil, errors.New("active kid cannot be blank")
	}
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}
	if lookup == nil {
		return nil, errors.New("public key function cannot be nil")
	}

	// Create the token (JWT) parser to use. The algorithm used to sign the JWT must be
	// validated to avoid a critical vulnerability:
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	// Authenticator
	a := Auth{
		privateKey: privateKey,
		activeKID:  activeKID,
		algorithm:  algorithm,
		keyFunc:    lookup,
		parser:     &parser,
	}

	return &a, nil
}

// GenerateToken generates a signed token (JWT) string representing the user Claims.
// **************************************************************************************
// TODO: Per KeyID; see `modelskit.Token(..)`; it's commented out; taken from service-5.
// **************************************************************************************
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)

	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.activeKID

	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}

// ValidateToken recreates the Claims from whence the token was created.
// It verifies that the token was signed using the key identified therein.
// If invalid, returns: `nil, err`; reason per error message.
func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {

	// `getPubKey` DEFINES a function that, WHEN INVOKED, RETURNs the PUBLIC KEY
	// per key ID (`kid`); invoking `PubKeyLookup` function, `keyFunc(ID)`.
	// The ID is extracted from the token's header section.
	getPubKey := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}
		kidStr, ok := kid.(string)
		if !ok {
			return nil, errors.New("token key id (kid) must be string")
		}
		return a.keyFunc(kidStr)
	}
	// VALIDATE the Claims by parsing token payload, `tokenStr`
	// USING the PUBLIC KEY recovered per `getPubKey` execution.
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, getPubKey)
	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}

// ----------------------------------------------------------------------------
// UNUSED

// Valid is called during the parsing of a token.
func (c Claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return errors.Wrap(err, "validating standard claims")
	}

	return nil
}
