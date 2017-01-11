package api

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/boxtown/meirl/data"
)

// PostAPI contains state information for executing
// MeIRL User API route handlers
type PostAPI struct {
	stores data.Stores
	debug  bool
}

// NewPostAPI returns an instance of the UserAPI struct
func NewPostAPI(stores data.Stores, debug bool) PostAPI {
	return PostAPI{
		stores: stores,
		debug:  debug,
	}
}

// CreatePost returns an http handler that handles create post API
// requests
func (api PostAPI) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := claimsID(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var p data.Post
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		p.AuthorID = userID

		id, err := api.stores.PostStore.Create(&p)
		if err != nil {
			writeError(err, w, api.debug)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Location", fmt.Sprintf("/%s/post/%d", apiVersion, id))
		writeJSON(IDResponse{ID: id}, w)
	}
}

// GetPost returns an http handler that handles get post API
// requests
func (api PostAPI) GetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := r.Context().Value(idContextKey).(int64)
		post, err := api.stores.PostStore.Get(id)
		if err == data.ErrNoEnt {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			writeError(err, w, api.debug)
			return
		}
		writeJSON(post, w)
	}
}
