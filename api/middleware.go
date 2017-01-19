package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// CORS is a middleware function that handles CORS logic
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement CORS
		next(w, r)
	}
}

// BodySizeLimiter is a middleware that limits the size of the body
// read from an http request
type BodySizeLimiter struct {
	h http.Handler
	n int64
}

func (bl BodySizeLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, bl.n)
	bl.h.ServeHTTP(w, r)
}

// LimitBodySize returns an http handler that wraps the given handler
// within a BodySizeLimiter to limit the size of the body read from the
// request to n bytes
func LimitBodySize(handler http.Handler, n int64) http.Handler {
	return &BodySizeLimiter{
		h: handler,
		n: n,
	}
}

// GetIDMiddleware is a middleware function that retrieves an ID from the
// route path and injects it into API functions that request an ID
func GetIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), idContextKey, id))
		next(w, r)
	}
}

// GetClaimsMiddleware is a middleware that attempts to parse a JWT from the
// 'Authorization' header and injects it into API functions requesting a
// JWT. Responds with a 400 Bad Request if the header is not found or
// invalid
func GetClaimsMiddleware(signingKey []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		parts := strings.Split(strings.TrimSpace(authHeader), " ")
		if parts[0] != "Bearer" || len(parts) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, ok := token.Claims.(jwt.MapClaims); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), claimsContextKey, token.Claims))
		next(w, r)
	}
}
