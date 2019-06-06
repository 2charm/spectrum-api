package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/2charm/spectrum-api/pkg/util"

	"github.com/2charm/spectrum-api/pkg/handlers"
	"github.com/2charm/spectrum-api/pkg/sessions"
	"github.com/2charm/spectrum-api/pkg/users"
	"github.com/go-redis/redis"
)

func main() {
	addr := util.GetEnvironmentVariable("ADDR")
	newsaddr := util.GetEnvironmentVariable("NEWSADDR")
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

	rs := sessions.NewRedisStore(rdb, time.Hour)

	//mySQL Server
	db, err := sql.Open("mysql", dsn)
	util.FailOnError(err, "Error opening a new SQL database")

	err = db.Ping()
	util.FailOnError(err, "Error pinging database")
	log.Printf("Successfully connected to SQL database!\n")

	ms := users.NewMySQLStore(db)

	ctx := handlers.HandlerContext{
		SigningKey:   sessionkey,
		SessionStore: rs,
		UserStore:    ms,
	}

	newsURL, err := url.Parse("http://" + newsaddr)
	util.FailOnError(err, "Invalid URL for microservice")

	log.Printf("News Microservice URL: %s", newsURL.String())
	newsProxy := &httputil.ReverseProxy{Director: customDirector(newsURL, &ctx)}

	mux := http.NewServeMux()
	mux.Handle("/v1/news", newsProxy)                           //Get news
	mux.Handle("/v1/spectrum/", newsProxy)                      //Get related news
	mux.Handle("/v1/metrics", newsProxy)                        //Get and post metrics
	mux.HandleFunc("/v1/users", ctx.UsersHandler)               //Create user
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)         //Login user
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler) //Logout user
	wrappedMux := handlers.NewResponseHeader(mux)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))
}

func customDirector(target *url.URL, ctx *handlers.HandlerContext) func(*http.Request) {
	return func(r *http.Request) {
		sessState := &handlers.SessionState{}
		log.Print("Accessing session...")
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessState)
		if err == nil {
			obj, err := json.Marshal(sessState.User)
			if err == nil {
				log.Print("Valid User!")
				r.Header.Del("X-User")
				r.Header.Add("X-User", string(obj))
			}
		} else {
			log.Printf("Error getting state: %v", err)
		}
		r.Host = target.Host
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
	}
}
