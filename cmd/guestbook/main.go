package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	rc *redis.Client
)

func HandleError(result interface{}, err error) (r interface{}) {
	if err != nil {
		panic(err)
	}
	return result
}

func ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v, err := rc.LRange(k, 0, -1).Result()
	if err != nil {
		HandleError(nil, err)
	}

	membersJSON := HandleError(json.MarshalIndent(v, "", "  ")).([]byte)
	rw.Write(membersJSON)
}

func ListPushHandler(rw http.ResponseWriter, req *http.Request) {
	k := mux.Vars(req)["key"]
	v := mux.Vars(req)["value"]
	_, err := rc.LPush(k, v).Result()
	if err != nil {
		HandleError(nil, err)
	}
	ListRangeHandler(rw, req)
}

func EnvHandler(rw http.ResponseWriter, req *http.Request) {
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

func main() {
	rc = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rc.Close()

	r := mux.NewRouter()
	r.Path("/env").Methods("GET").HandlerFunc(EnvHandler)
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(ListRangeHandler)
	r.Path("/lpush/{key}/{value}").Methods("GET").HandlerFunc(ListPushHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
