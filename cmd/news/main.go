package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"

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

	as := news.NewArticleStore(db)

	ctx := news.HandlerContext{
		APIKey:       apikey,
		ArticleStore: as,
		ArticleCache: cache.New(time.Minute*15, time.Minute*25),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/news", ctx.NewsHandler)          //Get news
	mux.HandleFunc("/v1/spectrum/", ctx.SpectrumHandler) //Get full spectrum of news
	mux.HandleFunc("/v1/metrics", ctx.MetricsHandler)    //Get and post metrics

	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
