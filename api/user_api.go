package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"regexp"

	"github.com/boxtown/meirl/data"
)

var errBadUsername = errors.New("Username may only contain [0-9], [a-z], and [A-Z]")
var errTakenUsername = errors.New("Username is taken")
var errBadEmail = errors.New("Email must be of the format [example@example]")
var errTakenEmail = errors.New("Email is taken")

// UserAPI contains state information for executing
// MeIRL User API route handlers
type UserAPI struct {
	stores data.Stores
	auth   Auth
	debug  bool
}

// NewUserAPI returns an instance of the UserAPI struct
func NewUserAPI(stores data.Stores, auth Auth, debug bool) UserAPI {
	return UserAPI{
		stores: stores,
		auth:   auth,
		debug:  debug,
	}
}

// CreateUser returns an http handler that handles create user API
// requests
func (api UserAPI) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u data.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = api.isCreateRequestValid(&u)
		if err != nil {
			if isCustomError(err) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, err.Error())
			} else {
				writeError(err, w, api.debug)
			}
			return
		}
		u.Password, err = api.auth.SecurePassword(u.Password)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		id, err := api.stores.UserStore.Create(&u)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Location", fmt.Sprintf("/%s/user/%d", apiVersion, id))
		writeJSON(IDResponse{ID: id}, w)
	}
}

// GetUser returns an http handler that handles get user API
// requests
func (api UserAPI) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := contextID(r)
		user, err := api.stores.UserStore.Get(id)
		if err == data.ErrNoEnt {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			writeError(err, w, api.debug)
			return
		}
		user.Password = ""
		writeJSON(user, w)
	}
}

// GetMe returns an http handler that uses JWT claims information
// to retrieve user information
func (api UserAPI) GetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := claimsID(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := api.stores.UserStore.Get(id)
		if err == data.ErrNoEnt {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			writeError(err, w, api.debug)
			return
		}
		user.Password = ""
		writeJSON(user, w)
	}
}

// GetFeed returns an http handler that handles get user feed
// API requests
func (api UserAPI) GetFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := contextID(r)
		_, err := api.stores.UserStore.Get(id)
		if err == data.ErrNoEnt {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			writeError(err, w, api.debug)
			return
		}
		options := ListOptionsFromRequest(r)
		posts, err := api.stores.PostStore.Feed(id, options, data.PostSortByDate)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		writeJSON(posts, w)
	}
}

// FollowUser returns an http handler that handles follower user
// API requests
func (api UserAPI) FollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := claimsID(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		followeeID := contextID(r)
		errc := make(chan error, 2)
		go func() {
			_, err := api.stores.UserStore.Get(id)
			errc <- err
		}()
		go func() {
			_, err := api.stores.UserStore.Get(followeeID)
			errc <- err
		}()
		for i := 0; i < 2; i++ {
			err := <-errc
			if err == data.ErrNoEnt {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		err := api.stores.Follow(id, followeeID)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

// UnFollowUser returns an http handler that handles unfollowing user
// API requests
func (api UserAPI) UnFollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteUser returns an http handler that handles delete user API
// requests
func (api UserAPI) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := contextID(r)
		err := api.stores.UserStore.Delete(id)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

// Login returns an http handler that handles user login API
// requests
func (api UserAPI) Login(signingKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u data.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		stored, err := api.getStoredUser(&u)
		if err != nil {
			if err == data.ErrNoEnt {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			writeError(err, w, api.debug)
			return
		}
		valid := api.auth.CheckPassword(u.Password, stored.Password)
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		accessToken, err := api.auth.GenerateAccessToken(stored, signingKey)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		writeJSON(TokenResponse{AccessToken: accessToken}, w)
	}
}

func (api UserAPI) getStoredUser(user *data.User) (*data.User, error) {
	if len(user.Username) > 0 {
		return api.stores.GetByUsername(user.Username)
	}
	return api.stores.GetByEmail(user.Email)
}

func (api UserAPI) isCreateRequestValid(user *data.User) error {
	usernameRegex := regexp.MustCompile("^[0-9a-zA-Z_]+$")
	emailRegex := regexp.MustCompile("^.+@.+$")
	if !usernameRegex.MatchString(user.Username) {
		return errBadUsername
	}
	if !emailRegex.MatchString(user.Email) {
		return errBadEmail
	}
	uc := make(chan error, 1)
	ec := make(chan error, 1)
	go func() {
		_, err := api.stores.GetByUsername(user.Username)
		uc <- err
	}()
	go func() {
		_, err := api.stores.GetByEmail(user.Email)
		ec <- err
	}()
	err := <-uc
	if err != data.ErrNoEnt {
		if err != nil {
			return err
		}
		return errTakenUsername
	}
	err = <-ec
	if err != data.ErrNoEnt {
		if err != nil {
			return err
		}
		return errTakenEmail
	}
	return nil
}

func isCustomError(err error) bool {
	return err == errBadUsername || err == errBadEmail || err == errTakenUsername || err == errTakenEmail
}
