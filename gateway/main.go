package main

import (
	"log"
	"net/http"
	"os"

	"github.com/2charm/spectrum-api/gateway/handlers"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getEnvironmentVariable(key string) string {
	val, set := os.LookupEnv(key)
	if !set {
		log.Fatalf("%s environment variable is not set", key)
	}
	return val
}

func main() {
	addr := getEnvironmentVariable("ADDR")
	key := getEnvironmentVariable("KEY")
	ctx := handlers.HandlerContext{key}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/news", ctx.NewsHandler)
	http.ListenAndServe(addr, mux)
}
