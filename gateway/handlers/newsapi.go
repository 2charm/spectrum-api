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
const topHeadlinesURL = baseURL + "top-headlines?apiKey=%s&language=en&pageSize=10&sources=google-news"
const everythingURL = baseURL + "everything?apiKey=%s&language=en&pageSize=5&q=%s"

//NewsHandler handles requests for the top headlines
func (ctx *HandlerContext) NewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		topHeadlines, err := getTopHeadlines(ctx.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{}
		max := 5
		for i := 0; i < max; i++ {
			title := topHeadlines.Articles[i].Title
			keywords, err := getKeywords(title)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if keywords != "" {
				log.Printf("Keywords: %s\n", keywords)
				articles, err := getRelatedArticles(keywords, ctx.Key)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				response[keywords] = articles
			} else {
				max++
			}
		}
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

func getTopHeadlines(key string) (*models.Headlines, error) {
	resp, err := http.Get(fmt.Sprintf(topHeadlinesURL, key))
	if err != nil {
		return nil, fmt.Errorf("Error calling NewsAPI Top headlines: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading json from NewsAPI: %v", err)
	}
	headlines := &models.Headlines{}
	err = json.Unmarshal(body, headlines)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling bytes: %v", err)
	}
	return headlines, nil
}

func getRelatedArticles(keywords string, key string) ([]models.Article, error) {
	reqURL := fmt.Sprintf(everythingURL, key, keywords)
	log.Print(reqURL)
	resp, err := http.Get(reqURL)
	log.Printf("%+v", resp)
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
