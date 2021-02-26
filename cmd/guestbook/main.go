package main

import (
	"fmt"
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
			Int64("duration", time.Since(start).Microseconds()).
			Send()
	})
}

func main() {
	if os.Getenv("ENV") == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	redisConn := os.Getenv("REDIS")
	if redisConn == "" {
		redisConn = "localhost:6379"
	}

	// redis_addr := os.Getenv("REDIS")
	rc := redis.NewClient(&redis.Options{
		Addr: redisConn,
		DB:   0,
	})
	defer rc.Close()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	rr := r.PathPrefix("/api").Subrouter()

	a, _ := api.New(rc)
	a.Init(rr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	log.Info().Msgf("Ready to serve on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil); err != nil {
		log.Fatal().Err(err)
	}
}
