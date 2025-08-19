package models

type NewsArticleResponse struct {
	Status            int             `json:"status"`
	StatusDescription string          `json:"statusDescription"`
	Data              NewsArticleData `json:"data"`
}

type NewsArticleData struct {
	Results  []NewsArticle `json:"results"`
	NextPage string        `json:"nextPage"`
}

type NewsArticle struct {
	ArticleID      string   `json:"article_id"`
	Title          string   `json:"title"`
	Link           string   `json:"link"`
	Keywords       []string `json:"keywords"`
	Creator        []string `json:"creator"`
	Description    string   `json:"description"`
	Content        string   `json:"content"`
	PubDate        string   `json:"pubDate"`
	PubDateTZ      string   `json:"pubDateTZ"`
	ImageURL       string   `json:"image_url"`
	VideoURL       *string  `json:"video_url"` // optional
	SourceID       string   `json:"source_id"`
	SourceName     string   `json:"source_name"`
	SourcePriority int      `json:"source_priority"`
	SourceURL      string   `json:"source_url"`
	SourceIcon     string   `json:"source_icon"`
	Language       string   `json:"language"`
	Country        []string `json:"country"`
	Category       []string `json:"category"`
	Sentiment      string   `json:"sentiment"`
	SentimentStats string   `json:"sentiment_stats"`
	AITag          string   `json:"ai_tag"`
	AIRegion       string   `json:"ai_region"`
	AIOrg          string   `json:"ai_org"`
	AISummary      string   `json:"ai_summary"`
	AIContent      string   `json:"ai_content"`
	Duplicate      bool     `json:"duplicate"`
}

type NewsDataResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Results      []NewsArticle `json:"results"`
	NextPage     string        `json:"nextPage"`
}
