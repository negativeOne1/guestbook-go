package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Service struct {
	rc *redis.Client
}

func New(rc *redis.Client) (*Service, error) {
	return &Service{rc: rc}, nil
}

func (s *Service) Init(r *mux.Router) {
	r.Path("/env").Methods("GET").HandlerFunc(s.EnvHandler)
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(s.ListRangeHandler)
	r.Path("/rpush/{key}/{value}").Methods("GET").HandlerFunc(s.ListPushHandler)
}

func (s *Service) ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v, err := s.rc.LRange(k, 0, -1).Result()
	if err != nil {
		log.Err(err).Msgf("could not get value for %s", k)
		return
	}

	membersJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Err(err).Msg("can't marshal json")
		return
	}

	if _, err := rw.Write([]byte(membersJSON)); err != nil {
		log.Err(err).Msg("can't write output")
		return
	}
}

func (s *Service) ListPushHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v := mux.Vars(req)["value"]
	_, err := s.rc.LPush(k, v).Result()
	if err != nil {
		log.Err(err).Msgf("can't get entry for k: %s, v: %s", k, v)
		return
	}

	s.ListRangeHandler(rw, req)
}

func (s *Service) EnvHandler(rw http.ResponseWriter, req *http.Request) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}
	envJSON, err := json.MarshalIndent(environment, "", "  ")
	if err != nil {
		log.Err(err).Msg("can't marshal json")
		return
	}

	if _, err := rw.Write([]byte(envJSON)); err != nil {
		log.Err(err).Msg("can't write output")
		return
	}
}
