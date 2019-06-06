package news

type Headlines struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source          source `json:"source"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	URLToImage      string `json:"urlToImage"`
	PublishedAt     string `json:"publishedAt"`
	Content         string `json:"content"`
	SpectrumEnabled bool   `json:"spectrumEnabled"`
}

type source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Metrics struct {
	UserID                int64          `json:"userID"`
	CategoryToNumArticles map[string]int `json:"categoryToNumArticles"`
	SourceToNumArticles   map[string]int `json:"sourceToNumArticles"`
}

type NewMetric struct {
	Category string `json:"category"`
	Source   string `json:"source"`
}
