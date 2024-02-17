package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/jwt"
	middleware "github.com/oapi-codegen/echo-middleware"
)

// JWSValidator is used to validate JWS payloads and return a JWT if they're
// valid
type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

const JWTClaimsContextKey = "jwt_claims"

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")
)

// GetJWSFromRequest retrieves the JWS from the Authorization header of an HTTP request.
// It expects the Authorization header value to be in the format "Bearer <token>",
// with a space following the "Bearer" keyword.
// If the Authorization header is missing, it returns an ErrNoAuthHeader error.
// If the Authorization header is malformed, it returns an ErrInvalidAuthHeader error.
// Otherwise, it trims the "Bearer " prefix from the header value and returns the JWS.
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

// NewAuthenticator creates a function that can be used as an authentication
// function in openapi3filter.Options. It takes a JWSValidator as input and
// returns an authentication function that calls Authenticate. The Authenticate
// function validates the security scheme name, gets the JWS from the request,
// validates the JWS, checks the token claims against the expected claims, sets
// the JWT claims on the request context, and returns an error if any of these
// steps fail.
func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	// We've got a valid token now, and we can look into its claims to see whether
	// they match. Every single scope must be present in the claims.
	err = CheckTokenClaims(input.Scopes, token)

	if err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}

	// Set the property on the echo context so the handler is able to
	// access the claims data we generate in here.
	eCtx := middleware.GetEchoContext(ctx)
	eCtx.Set(JWTClaimsContextKey, token)

	return nil
}

// GetClaimsFromToken returns a list of claims from the token. We store these
// as a list under the "perms" claim, short for permissions, to keep the token
// shorter.
func GetClaimsFromToken(t jwt.Token) ([]string, error) {
	rawPerms, found := t.Get(PermissionsClaim)
	if !found {
		// If the perms aren't found, it means that the token has none, but it has
		// passed signature validation by now, so it's a valid token, so we return
		// the empty list.
		return make([]string, 0), nil
	}

	// rawPerms will be an untyped JSON list, so we need to convert it to
	// a string list.
	rawList, ok := rawPerms.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' claim is unexpected type'", PermissionsClaim)
	}

	claims := make([]string, len(rawList))

	for i, rawClaim := range rawList {
		var ok bool
		claims[i], ok = rawClaim.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] is not a string", PermissionsClaim, i)
		}
	}
	return claims, nil
}

// CheckTokenClaims verifies that the provided JWT token contains all the expected claims.
// It first extracts the claims from the token using GetClaimsFromToken function. If the
// extraction fails, it returns an error detailing the failure to get claims. Then, it
// checks whether each of the expected claims is present in the token's claims. If any of
// the expected claims are missing, it returns an ErrClaimsInvalid error indicating that
// the token does not have the required claims.
func CheckTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := GetClaimsFromToken(t)
	if err != nil {
		return fmt.Errorf("getting claims from token: %w", err)
	}
	// Put the claims into a map, for quick access.
	claimsMap := make(map[string]bool, len(claims))
	for _, c := range claims {
		claimsMap[c] = true
	}

	for _, e := range expectedClaims {
		if !claimsMap[e] {
			return ErrClaimsInvalid
		}
	}
	return nil
}
