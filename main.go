package main

import (
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"

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
	graceful.Run(":8080", 10*time.Second, cors.Default().Handler(r))
}
