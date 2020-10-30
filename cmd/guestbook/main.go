package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/negativeOne1/guestbook-go/pkg/api"
)

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rc.Close()

	s, _ := api.New(rc)
	r := mux.NewRouter()
	r.Path("/env").Methods("GET").HandlerFunc(s.EnvHandler)
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(s.ListRangeHandler)
	r.Path("/lpush/{key}/{value}").Methods("GET").HandlerFunc(s.ListPushHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
