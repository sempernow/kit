package auth_test

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"gd9/prj3/app"
	"gd9/prj3/kit/auth"
	"gd9/prj3/kit/convert"
	"gd9/prj3/kit/testkit"

	"github.com/golang-jwt/jwt"
)

func TestAuthenticator(t *testing.T) {
	t.Skip()
	//*******************************************************************************
	// REQUIREs: export APP_AUTH_PRIVATE_KEY_FILE=$(cat './assets/keys/private.pem')
	//*******************************************************************************
	var (
		privateKey *rsa.PrivateKey
		bb         []byte
		err        error
	)

	bb, err = ioutil.ReadFile(os.Getenv("APP_AUTH_PRIVATE_KEY_FILE"))
	if err != nil {
		fmt.Println("... is STRING")
		bb = []byte(os.Getenv("APP_AUTH_PRIVATE_KEY_FILE"))
	} else {
		fmt.Println("... is FILE")
	}
	fmt.Println(convert.BytesToString(bb))
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bb)
	if err != nil {
		t.Fatal("FAIL @ privateKey")
	}
	//t.Fatal("=== END")

	//privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pvtPEM))
	testkit.Log(t, "Parse PEM-encoded private key", err)

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
	testkit.Log(t, "Parse PEM-encoded public key", err)

	keyLookupFunc := func(kid string) (*rsa.PublicKey, error) {
		if kid != KID {
			return nil, errors.New("no public key found")
		}
		return publicKey, nil
	}
	a, err := auth.New(privateKey, KID, "RS256", keyLookupFunc)
	testkit.Log(t, "Create Authenticator", err)

	t.Log("@ Authenticate and authorize user access.")
	{
		t.Logf("\t@ Valid claims")
		{
			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    auth.Issuer(auth.PWA),
					Subject:   "0x01",
					Audience:  app.Audience,
					ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(), // One year.
					IssuedAt:  auth.IssuedAtNow(),
				},
				Roles: []string{auth.RoleAppOps},
			}
			//claims = auth.Claims{StandardClaims: jwt.StandardClaims{}, Roles: []string{}}
			//... panic ... index out of range ...
			token, err := a.GenerateToken(claims)
			testkit.Log(t, "@ Generate the JWT", err)

			parsedClaims, err := a.ValidateToken(token)
			testkit.Log(t, "@ Validate/Parse claims", err)

			msg := "Want expected number of user Roles"
			{
				exp, got := len(claims.Roles), len(parsedClaims.Roles)
				testkit.LogDiff(t, msg, got, exp)
			}
			msg = "Want expected user Roles"
			{
				exp, got := claims.Roles[0], parsedClaims.Roles[0]
				testkit.LogDiff(t, msg, got, exp)
			}

		}
	}
	t.Logf("\t@ Invalid (expired) claims")
	{
		claims := auth.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    auth.Issuer(auth.PWA),
				Subject:   "0x01",
				Audience:  app.Audience,
				ExpiresAt: time.Now().Add(-1 * time.Minute).Unix(), // @ Expired
				IssuedAt:  auth.IssuedAtNow(),
			},
			Roles: []string{auth.RoleAppOps},
		}
		//claims = auth.Claims{StandardClaims: jwt.StandardClaims{}, Roles: []string{}}
		//... panic ... index out of range ...
		token, err := a.GenerateToken(claims)
		testkit.Log(t, "@ Generate the JWT", err)

		_, err = a.ValidateToken(token)
		msg := "@ Validate/Parse claims : Want '... token is expired ...'"
		testkit.LogDiff(t, msg, strings.Contains(err.Error(), "token is expired"), true)
	}
	t.Logf("\t@ Nonexistent claims")
	{
		claims := auth.Claims{}
		//claims = auth.Claims{StandardClaims: jwt.StandardClaims{}, Roles: []string{}}
		token, err := a.GenerateToken(claims)
		testkit.Log(t, "Generate the JWT", err)

		_, err = a.ValidateToken(token)
		msg := "@ Validate/Parse claims : Want <nil>"
		testkit.LogDiff(t, msg, nil, err)
	}
}

// The key id we would have generated for the keys below.
const KID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

// Private PEM : Output of:
// openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
// ./admin keygen
var pvtPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAvMAHb0IoLvoYuW2kA+LTmnk+hfnBq1eYIh4CT/rMPCxgtzjq
U0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiMQl/bxAhgLkGrUhgpZ10Bt6GzSEgw
QNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNcLxPUmfjbbH64g+seg28pbgQPwKFK
tF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R42LH9c75JsRzz4FeN+VzwMAL6yQn
ZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLRBzQGN00Z7TuBRIfDJWU21qQN4Cq7
keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMGvwIDAQABAoIBAFQmQKpHkmavNYql
6POaksBRwaA1YzSijr7XJizGIXvKRSwqgb2zdnuTSgpspAx09Dr/aDdy7rZ0DAJt
fk2mInINDottOIQm3txwzTS58GQQAT/+fxTKWJMqwPfxYFPWqbbU76T8kXYna0Gs
OcK36GdMrgIfQqQyMs0Na8MpMg1LmkAxuqnFCXS/NMyKl9jInaaTS+Kz+BSzUMGQ
zebfLFsf2N7sLZuimt9zlRG30JJTfBlB04xsYMo734usA2ITe8U0XqG6Og0qc6ev
6lsoM8hpvEUsQLcjQQ5up7xx3S2stZJ8o0X8GEX5qUMaomil8mZ7X5xOlEqf7p+v
lXQ46cECgYEA2lbZQON6l3ZV9PCn9j1rEGaXio3SrAdTyWK3D1HF+/lEjClhMkfC
XrECOZYj+fiI9n+YpSog+tTDF7FTLf7VP21d2gnhQN6KAXUnLIypzXxodcC6h+8M
ZGJh/EydLvC7nPNoaXx96bohxzS8hrOlOlkCbr+8gPYKf8qkbe7HyxECgYEA3U6e
x9g4FfTvI5MGrhp2BIzoRSn7HlNQzjJ71iMHmM2kBm7TsER8Co1PmPDrP8K/UyGU
Q25usTsPSrHtKQEV6EsWKaP/6p2Q82sDkT9bZlV+OjRvOfpdO5rP6Q95vUmMGWJ/
S6oimbXXL8p3gDafw3vC1PCAhoaxMnGyKuZwlM8CgYEAixT1sXr2dZMg8DV4mMfI
8pqXf+AVyhWkzsz+FVkeyAKiIrKdQp0peI5C/5HfevVRscvX3aY3efCcEfSYKt2A
07WEKkdO4LahrIoHGT7FT6snE5NgfwTMnQl6p2/aVLNun20CHuf5gTBbIf069odr
Af7/KLMkjfWs/HiGQ6zuQjECgYEAv+DIvlDz3+Wr6dYyNoXuyWc6g60wc0ydhQo0
YKeikJPLoWA53lyih6uZ1escrP23UOaOXCDFjJi+W28FR0YProZbwuLUoqDW6pZg
U3DxWDrL5L9NqKEwcNt7ZIDsdnfsJp5F7F6o/UiyOFd9YQb7YkxN0r5rUTg7Lpdx
eMyv0/UCgYEAhX9MPzmTO4+N8naGFof1o8YP97pZj0HkEvM0hTaeAQFKJiwX5ijQ
xumKGh//G0AYsjqP02ItzOm2mWnbI3FrNlKmGFvR6VxIZMOyXvpLofHucjJ5SWli
eYjPklKcXaMftt1FVO4n+EKj1k1+Tv14nytq/J5WN+r4FBlNEYj/6vg=
-----END RSA PRIVATE KEY-----`

// Public PEM : Output of:
// openssl rsa -pubout -in private.pem -out public.pem
// ./admin keygen
const pubPEM = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvMAHb0IoLvoYuW2kA+LT
mnk+hfnBq1eYIh4CT/rMPCxgtzjqU0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiM
Ql/bxAhgLkGrUhgpZ10Bt6GzSEgwQNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNc
LxPUmfjbbH64g+seg28pbgQPwKFKtF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R
42LH9c75JsRzz4FeN+VzwMAL6yQnZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLR
BzQGN00Z7TuBRIfDJWU21qQN4Cq7keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMG
vwIDAQAB
-----END RSA PUBLIC KEY-----`
