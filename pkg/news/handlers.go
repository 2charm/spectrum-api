package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"

	"github.com/2charm/spectrum-api/pkg/users"
	prose "gopkg.in/jdkato/prose.v2"
)

const baseURL = "https://newsapi.org/v2/"
const cacheKey = "articles"

var categories = []string{"sports", "health", "business", "entertainment", "science", "technology"} //todo: add US and WORLD

//HandlerContext provides context for news handler package
type HandlerContext struct {
	APIKey       string
	ArticleStore Store
	ArticleCache *cache.Cache
}

func getUserFromHeader(r *http.Request) (*users.User, error) {
	val := r.Header.Get("X-User")
	if len(val) == 0 {
		log.Print("No X-User in Header")
		return nil, fmt.Errorf("No X-User in Header")
	}
	user := users.User{}
	err := json.Unmarshal([]byte(val), &user)
	if err != nil {
		log.Printf("Error unmarshalling X-User")
		return nil, fmt.Errorf("Error unmarshalling X-User")
	}
	return &user, nil
}

//NewsHandler handles requests for the articles needed by client
func (ctx *HandlerContext) NewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
		return
	}
	log.Print("GET /v1/news")

	var response map[string]interface{}
	//cachedArticles, expire, exists := ctx.ArticleCache.GetWithExpiration(cacheKey)
	if cachedArticles, exists := ctx.ArticleCache.Get(cacheKey); exists {
		response = cachedArticles.(map[string]interface{})
	} else {
		response = map[string]interface{}{}
		for _, category := range categories {
			articles, err := getArticlesByCategory(category, ctx.APIKey)
			if err != nil {
				log.Printf("API call went wrong for %s category: %v", category, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			articles = checkSpectrumEnabled(articles)
			response[category] = articles
		}
		articles, err := callNewsAPI("top-headlines", "country=us&category=general&pageSize=10&apiKey="+ctx.APIKey)
		if err != nil {
			log.Printf("API call went wrong for headlines: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles = checkSpectrumEnabled(articles)
		response["headline"] = articles

		articles, err = getArticlesByCategory("general", ctx.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles = checkSpectrumEnabled(articles)
		response["us"] = articles
		err = ctx.ArticleCache.Add(cacheKey, response, time.Hour*3)
		if err != nil {
			log.Print("Error inserting articles to cache")
		}
	}
	buffer, err := json.Marshal(response)
	if err != nil {
		log.Print("Marshal error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buffer)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

}

//MetricsHandler handles requests for metrics by users
func (ctx *HandlerContext) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromHeader(r)
	if err != nil {
		log.Print("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	if r.Method == "POST" {
		log.Print("POST /v1/metrics")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		metric := NewMetric{}
		err = json.Unmarshal(body, &metric)
		if err != nil {
			log.Printf("Error unmarshalling json metric: %v", err)
			http.Error(w, "Error unmarshalling json", http.StatusBadRequest)
			return
		}

		err = ctx.ArticleStore.InsertArticle(metric, user.ID)
		if err != nil {
			log.Printf("Error inserting article: %v", err)
			http.Error(w, "can't insert article", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == "GET" {
		metrics, err := ctx.ArticleStore.GetByUserID(user.ID)
		if err != nil {
			log.Printf("Error retrieving metrics: %v", err)
			http.Error(w, "can't retrieve metrics", http.StatusInternalServerError)
			return
		}
		buffer, err := json.Marshal(metrics)
		if err != nil {
			log.Print("Marshal error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(buffer)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	} else {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
		return
	}
}

//SpectrumHandler handles requests for related articles needed by client
func (ctx *HandlerContext) SpectrumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid http method.", http.StatusMethodNotAllowed)
		return
	}
	title := path.Base(r.URL.String())

	var response []Article
	var err error
	if cachedArticles, exists := ctx.ArticleCache.Get(title); exists {
		response = cachedArticles.([]Article)
	} else {
		response, err = getRelatedArticles(title, ctx.APIKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving related articles:%s", err.Error()), http.StatusInternalServerError)
			return
		}
		ctx.ArticleCache.Add(title, response, time.Hour*15)
	}

	buffer, err := json.Marshal(response)
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
	reqURL := baseURL + endpoint + "?" + query
	resp, err := http.Get(reqURL)
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
	title = strings.Replace(title, "%20", " ", -1)
	title = strings.Replace(title, "'", " ", -1)
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

func getRelatedArticles(title string, key string) ([]Article, error) {
	keywords, err := getKeywords(title)
	log.Printf("Related keywords to %s: %s", title, keywords)
	if err != nil {
		return nil, err
	}
	return callNewsAPI("everything", "sortBy=relevancy&language=en&pageSize=10&q="+keywords+"&apiKey="+key)
}

func checkSpectrumEnabled(articles []Article) []Article {
	for i, article := range articles {
		keywords, err := getKeywords(article.Title[:strings.LastIndex(article.Title, "-")])
		if err == nil && strings.Contains(keywords, "%20") {
			log.Print(keywords, " enabled")
			articles[i].SpectrumEnabled = true
		}
	}
	return articles
}
