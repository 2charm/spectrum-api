package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	prose "gopkg.in/jdkato/prose.v2"
)

const baseURL = "https://newsapi.org/v2/"

var categories = []string{"sports", "health", "business", "entertainment", "science", "technology"} //todo: add US and WORLD

//HandlerContext provides context for news handler package
type HandlerContext struct {
	APIKey       string
	ArticleStore Store
}

//NewsHandler handles requests for the articles needed by client
func (ctx *HandlerContext) NewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := map[string]interface{}{}
		for _, category := range categories {
			articles, err := getArticlesByCategory(category, ctx.APIKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response[category] = articles
		}
		articles, err := callNewsAPI("top-headlines", "country=us&category=general&pageSize=10&apiKey="+ctx.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response["headline"] = articles

		articles, err = getArticlesByCategory("general", ctx.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response["us"] = articles

		buffer, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(buffer)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	} else if r.Method == "POST" {

	} else {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
	}
}

//SpectrumHandler handles requests for related articles needed by client
func (ctx *HandlerContext) SpectrumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
		return
	}
	article := Article{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&article); err != nil {
		http.Error(w, "Error decoding json into Article.", http.StatusBadRequest)
		return
	}
	articles, err := getRelatedArticles(article, ctx.APIKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving related articles:%s", err.Error()), http.StatusInternalServerError)
		return
	}
	buffer, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buffer)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func getArticlesByCategory(category string, key string) ([]Article, error) {
	return callNewsAPI("top-headlines", "country=us&pageSize=10&category="+category+"&apiKey="+key)
}

func callNewsAPI(endpoint string, query string) ([]Article, error) {
	resp, err := http.Get(baseURL + endpoint + "?" + query)
	if err != nil {
		return nil, fmt.Errorf("Error calling NewsAPI everything: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading json from NewsAPI: %v", err)
	}
	headlines := &Headlines{}
	err = json.Unmarshal(body, headlines)
	if err != nil {
		log.Print(string(body))
		return nil, fmt.Errorf("Error unmarshalling bytes: %v", err)
	}
	return headlines.Articles, nil
}

func getKeywords(title string) (string, error) {
	doc, err := prose.NewDocument(title)
	if err != nil {
		return "", fmt.Errorf("Error retrieving NLTP: %v", err)
	}
	keywords := []string{}
	for _, word := range doc.Entities() {
		word.Text = strings.Replace(word.Text, " ", "%20", -1)
		keywords = append(keywords, word.Text)
	}
	return strings.Join(keywords, "%20"), nil
}

func getRelatedArticles(article Article, key string) ([]Article, error) {
	keywords, err := getKeywords(article.Title)
	if err != nil {
		return nil, err
	}
	return callNewsAPI("top-headlines", "pageSize=10&q="+keywords+"&apiKey="+key)
}
