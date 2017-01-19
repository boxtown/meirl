package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/boxtown/meirl/data"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/uber-go/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	apiVersion = "1.0"
)

var logger = zap.New(zap.NewTextEncoder())

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

// Attempt to write the body as JSON to the response writer.
// If WriteHeader has not been called, a 200 status will be auto-sent
func writeJSON(body interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

// Write a 503 error response to the response writer. If debug is true,
// will write the error message as well
func writeError(err error, w http.ResponseWriter, debug bool) {
	w.WriteHeader(http.StatusServiceUnavailable)
	if debug {
		logger.Error(err.Error())
		fmt.Fprintf(w, err.Error())
	}
}
