package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/2charm/spectrum-api/gateway/models"
	prose "gopkg.in/jdkato/prose.v2"
)

const baseURL = "https://newsapi.org/v2/"

var categories = []string{"sports", "health", "business", "entertainment", "science", "technology"} //todo: add US and WORLD

//NewsHandler handles requests for the categories needed by client
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
		articles, err := callNewsAPI("top-headlines", "language=en&pageSize=10&apiKey="+ctx.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response["headline"] = articles

		articles, err = callNewsAPI("top-headlines", "country=us&pageSize=10&category=general&apiKey="+ctx.APIKey)
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
	} else {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
	}
}

func getArticlesByCategory(category string, key string) ([]models.Article, error) {
	return callNewsAPI("top-headlines", "language=en&pageSize=10&category="+category+"&apiKey="+key)
}

func getArticlesByKeyword(keyword string, key string) ([]models.Article, error) {
	return callNewsAPI("everything", "language=en&q=seattle&sortBy=relevancy")
}

func callNewsAPI(endpoint string, query string) ([]models.Article, error) {
	resp, err := http.Get(baseURL + endpoint + "?" + query)
	if err != nil {
		return nil, fmt.Errorf("Error calling NewsAPI everything: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading json from NewsAPI: %v", err)
	}
	headlines := &models.Headlines{}
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

func getRelatedArticles(keywords string, key string) ([]models.Article, error) {
	reqURL := fmt.Sprintf(everythingURL, key, keywords)
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("Error calling NewsAPI everything: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading json from NewsAPI: %v", err)
	}
	headlines := &models.Headlines{}
	err = json.Unmarshal(body, headlines)
	if err != nil {
		log.Print(string(body))
		return nil, fmt.Errorf("Error unmarshalling bytes: %v", err)
	}
	return headlines.Articles, nil
}
