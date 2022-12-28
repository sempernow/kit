package auth

// TODO: Renew JWT  https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
// Create a '/refresh' route that takes the previous (valid) token
// and returns a new token with a renewed expiry time.
//
// To minimize misuse of a JWT, expiry time is kept to less than an hour or so.
// Typically the client application would refresh the token in the background.

// NOTE, unlike token (JWT) AUTHENTICATION,
//   token VALIDATION references nothing outside the JWT itself.
//   So, for example, shutting down the server, wiping the users database,
//   and restarting anew does NOT invalidate the (non-existent user's) token.
//   Only its own claims (ExpiresAt) can invalidate it (per token).
//
//   Deleting its private (paired) key and restarting this (key cacheing) server
//   de-authenticates that previously generated token, so it will fail validation here,
//   but such a "nuclear option" does so of ALL tokens created with that private key.
//
//   These dynamics reveal why the `ExpiresAt` claim is so vital to any JWT-based security.
//   Its time value, along with key-rotation schedule, determine the scheme's actual security.
