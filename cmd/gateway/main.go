package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/2charm/spectrum-api/pkg/util"

	"github.com/2charm/spectrum-api/pkg/handlers"
	"github.com/2charm/spectrum-api/pkg/sessions"
	"github.com/2charm/spectrum-api/pkg/users"
	"github.com/go-redis/redis"
)

func main() {
	addr := util.GetEnvironmentVariable("ADDR")
	apikey := util.GetEnvironmentVariable("APIKEY")
	tlscert := util.GetEnvironmentVariable("TLSCERT")
	tlskey := util.GetEnvironmentVariable("TLSKEY")
	sessionkey := util.GetEnvironmentVariable("SESSIONKEY")
	redisaddr := util.GetEnvironmentVariable("REDISADDR")
	dsn := util.GetEnvironmentVariable("DSN")

	//Redis Server
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	util.FailOnError(err, "Error pinging redis database")

	rs := sessions.NewRedisStore(rdb, 150*time.Second)

	//mySQL Server
	db, err := sql.Open("mysql", dsn)
	util.FailOnError(err, "Error opening a new SQL database")

	err = db.Ping()
	util.FailOnError(err, "Error pinging database")
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
