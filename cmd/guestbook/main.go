package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/negativeOne1/guestbook-go/pkg/api"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.EscapedPath(), time.Since(start))
	})
}

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rc.Close()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	rr := r.PathPrefix("/api").Subrouter()

	s, _ := api.New(rc)
	s.Init(rr)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
