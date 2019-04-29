package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/2charm/spectrum-api/gateway/handlers"
	"github.com/2charm/spectrum-api/gateway/models/sessions"
	"github.com/2charm/spectrum-api/gateway/models/users"
	"github.com/go-redis/redis"
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
	apikey := getEnvironmentVariable("APIKEY")
	tlscert := getEnvironmentVariable("TLSCERT")
	tlskey := getEnvironmentVariable("TLSKEY")
	sessionkey := getEnvironmentVariable("SESSIONKEY")
	redisaddr := getEnvironmentVariable("REDISADDR")
	dsn := getEnvironmentVariable("DSN")

	//Redis Server
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	failOnError(err, "Error pinging redis database")

	rs := sessions.NewRedisStore(rdb, 150*time.Second)

	//mySQL Server
	db, err := sql.Open("mysql", dsn)
	failOnError(err, "Error opening a new SQL database")

	err = db.Ping()
	failOnError(err, "Error pinging database")
	log.Printf("Successfully connected to SQL database!\n")

	ms := users.NewMySQLStore(db)

	ctx := handlers.HandlerContext{
		APIKey:       apikey,
		SigningKey:   sessionkey,
		SessionStore: rs,
		UserStore:    ms,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/news", ctx.NewsHandler)                //Get news
	mux.HandleFunc("/v1/users", ctx.UsersHandler)              //Create user
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)        //Login user
	mux.HandleFunc("v1/sessions/", ctx.SpecificSessionHandler) //Logout user
	wrappedMux := handlers.NewResponseHeader(mux)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))
}
