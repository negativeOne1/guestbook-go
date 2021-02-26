package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/negativeOne1/guestbook-go/pkg/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		next.ServeHTTP(w, r)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.EscapedPath()).
			Int64("duration", int64(time.Since(start))).
			Send()
	})
}

func main() {
	if os.Getenv("ENV") == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	// redis_addr := os.Getenv("REDIS")
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer rc.Close()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	rr := r.PathPrefix("/api").Subrouter()

	a, _ := api.New(rc)
	a.Init(rr)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	log.Info().Msg("Ready to serve on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal().Err(err)
	}
}
