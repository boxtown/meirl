package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"strings"

	"context"

	"github.com/boxtown/meirl/data"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	apiVersion = "1.0"
)

// PrefixAPIPath appends API-specific prefixes to the
// given path
func PrefixAPIPath(path string) string {
	return fmt.Sprintf("/%s/%s", apiVersion, path)
}

// ListOptionsFromRequest parses list options from
// a HTTP request. If the marker param is present, it is
// stored as a string and must be converted by the caller.
func ListOptionsFromRequest(request *http.Request) data.ListOptions {
	values := request.URL.Query()
	offset, err := strconv.Atoi(values.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	desc, _ := strconv.ParseBool(values.Get("desc"))
	options := data.ListOptions{
		Offset: offset,
		Limit:  limit,
		Desc:   desc,
		Marker: values.Get("marker"),
	}
	return options
}

// Auth is an interface for API authentication
type Auth interface {
	SecurePassword(password string) (string, error)
	CheckPassword(password, storedPassword string) bool
	GenerateAccessToken(user *data.User, signingKey []byte) (string, error)
}

// NewAuth returns a default implementation of Auth
// for the API
func NewAuth() Auth {
	return authImpl{}
}

// authImpl is the default API Auth implementation
type authImpl struct{}

// SecurePassword secures a password by performing a one-way
// hash on the password using BCrypt. Returns an error if there was
// an issue hashing the password
func (auth authImpl) SecurePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword returns if a given password matches the
// stored password hash
func (auth authImpl) CheckPassword(password, storedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) == nil
}

// GenerateAccessToken generates a JWT for the given user using
// the given signing key. Returns an error if there was an issue
// generating the JWT
func (auth authImpl) GenerateAccessToken(user *data.User, signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"iss": "MeIRL API Server",
	})
	return token.SignedString(signingKey)
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
		r = r.WithContext(context.WithValue(r.Context(), claimsContextKey, token.Claims))
		next(w, r)
	}
}

type contextKey int64

const (
	idContextKey contextKey = iota
	claimsContextKey
)

var errBadContext = errors.New("Could not retrieve value from context")

func writeJSON(body interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func writeError(err error, w http.ResponseWriter, debug bool) {
	w.WriteHeader(http.StatusServiceUnavailable)
	if debug {
		fmt.Fprintf(w, err.Error())
	}
}
