package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type Service struct {
	rc *redis.Client
}

func New(rc *redis.Client) (*Service, error) {
	return &Service{rc: rc}, nil
}

func (s *Service) Init(r *mux.Router) error {
	r.Path("/env").Methods("GET").HandlerFunc(s.EnvHandler)
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(s.ListRangeHandler)
	r.Path("/rpush/{key}/{value}").Methods("GET").HandlerFunc(s.ListPushHandler)

	return nil
}

func HandleError(result interface{}, err error) (r interface{}) {
	if err != nil {
		panic(err)
	}
	return result
}

func (s *Service) ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v, err := s.rc.LRange(k, 0, -1).Result()
	if err != nil {
		HandleError(nil, err)
	}

	membersJSON := HandleError(json.MarshalIndent(v, "", "  ")).([]byte)
	rw.Write(membersJSON)
}

func (s *Service) ListPushHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v := mux.Vars(req)["value"]
	_, err := s.rc.LPush(k, v).Result()
	if err != nil {
		HandleError(nil, err)
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

	envJSON := HandleError(json.MarshalIndent(environment, "", "  ")).([]byte)
	rw.Write(envJSON)
}
