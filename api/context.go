package api

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type contextKey int64

const (
	idContextKey contextKey = iota
	claimsContextKey
)

// Retrieve an ID from the context. Will panic if there was
// no ID stored using idContextKey or if the stored ID is not
// an int64
func contextID(r *http.Request) int64 {
	return r.Context().Value(idContextKey).(int64)
}

// Retrieve the user ID stored within the claims within an http
// context. Will panic if there are no claims stored within the context.
// Returns false if there is no ID witin the claims or if the ID is not
// an int64
func claimsID(r *http.Request) (int64, bool) {
	claims := r.Context().Value(claimsContextKey).(jwt.MapClaims)
	id, ok := claims["sub"].(int64)
	return id, ok
}
