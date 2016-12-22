package apitest

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// RequestWithContextID returns an http.Request from the given request with the id stored at key within
// the request's context
func RequestWithContextID(r *http.Request, key interface{}, id int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, id))
}

// RequestWithClaims returns an http.Request from the given request with the claims stored at key within
// the request's context
func RequestWithClaims(r *http.Request, key interface{}, claims jwt.MapClaims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, claims))
}
