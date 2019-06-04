package main

import (
	"log"
	"net/http"

	"github.com/2charm/spectrum-api/pkg/news"
	"github.com/2charm/spectrum-api/pkg/util"
)

func main() {
	addr := util.GetEnvironmentVariable("ADDR")
	apikey := util.GetEnvironmentVariable("APIKEY")
	dsn := util.GetEnvironmentVariable("DSN")

	//mySQL Server
	db, err := sql.Open("mysql", dsn)
	util.FailOnError(err, "Error opening a new SQL database")

	err = db.Ping()
	util.FailOnError(err, "Error pinging database")
	log.Printf("Successfully connected to SQL database!\n")

	ctx := news.HandlerContext{
		APIKey: apikey,
		NewsStore: 
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/news", ctx.NewsHandler)         //Get news
	mux.HandleFunc("/v1/spectrum", ctx.SpectrumHandler) //Get full spectrum of news

	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
