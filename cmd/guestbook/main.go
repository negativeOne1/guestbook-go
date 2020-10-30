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

	r := mux.NewRouter()
	s, _ := api.New(rc)
	s.Init(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
