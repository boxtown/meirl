package main

import (
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"

	"net/http"

	"github.com/boxtown/meirl/data"
	"github.com/boxtown/meirl/data/postgres"
	"github.com/rs/cors"
)

func main() {
	db, err := postgres.InitDB(pgUser, pgPass, pgHost, pgPort, pgDBName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userStore := postgres.NewUserStore(db)
	postStore := postgres.NewPostStore(db)
	r := Router(data.Stores{
		UserStore: userStore,
		PostStore: postStore,
	})
	graceful.Run(":8080", 10*time.Second, limitBodySize(cors.Default().Handler(r), requestBodyMaxBytes))
}

type bodyLimiter struct {
	h http.Handler
	n int64
}

func (bl bodyLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, bl.n)
	bl.h.ServeHTTP(w, r)
}

func limitBodySize(handler http.Handler, n int64) http.Handler {
	return &bodyLimiter{
		h: handler,
		n: n,
	}
}
