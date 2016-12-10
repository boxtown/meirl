package main

import (
	"net/http"

	"github.com/boxtown/meirl/api"
	"github.com/boxtown/meirl/data"
	"github.com/gorilla/mux"
)

// Router initializes an http.Handler that routes requests to the proper
// MeIRL request handlers
func Router(stores data.Stores) http.Handler {
	r := mux.NewRouter()
	initUserRoutes(r, stores)
	initPostRoutes(r, stores)
	return r
}

func initUserRoutes(r *mux.Router, stores data.Stores) {
	userAPI := api.NewUserAPI(stores, api.NewAuth(), debug())
	r.HandleFunc(
		api.PrefixAPIPath("user/{id:[0-9]+}"),
		api.GetIDMiddleware(userAPI.GetUser())).Methods("GET")
	r.HandleFunc(
		api.PrefixAPIPath("user/me"),
		api.GetClaimsMiddleware(signingKey, userAPI.GetMe("/user/"))).Methods("GET")
	r.HandleFunc(
		api.PrefixAPIPath("user/{id:[0-9]+}/feed"),
		api.GetIDMiddleware(userAPI.GetFeed())).Methods("GET")
	r.HandleFunc(
		api.PrefixAPIPath("user/new"),
		userAPI.CreateUser()).Methods("POST")
	r.HandleFunc(
		api.PrefixAPIPath("user/login"),
		userAPI.Login(signingKey)).Methods("POST")
	r.HandleFunc(
		api.PrefixAPIPath("user/{id:[0-9]+}"),
		api.GetIDMiddleware(userAPI.DeleteUser())).Methods("DELETE")
}

func initPostRoutes(r *mux.Router, stores data.Stores) {
	postAPI := api.NewPostAPI(stores, debug())
	r.HandleFunc(
		api.PrefixAPIPath("post/{id:[0-9]+}"),
		api.GetIDMiddleware(postAPI.GetPost())).Methods("GET")
	r.HandleFunc(
		api.PrefixAPIPath("post/new"),
		api.GetClaimsMiddleware(signingKey, postAPI.CreatePost())).Methods("POST")
}
