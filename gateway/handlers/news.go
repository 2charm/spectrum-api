package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const newsAPIURL = "https://newsapi.org/v2"

//NewsHandler handles requests for the top headlines
func (ctx *HandlerContext) NewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		resp, err := http.Get(fmt.Sprintf("%s/top-headlines?pageSize=%s&apiKey=%s", newsAPIURL, "10", ctx.Key))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error calling NewsAPI: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading json from NewsAPI: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write(body)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		return
	} else {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
		return
	}
}
