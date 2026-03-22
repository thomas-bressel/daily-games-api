package pkg

import "time"

// Article represents a parsed RSS article ready to be served by the API.
type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pubDate"`
	Creator     string    `json:"creator"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	Source      string    `json:"source"`
	SourceName  string    `json:"sourceName"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
}

// Feed represents a configured RSS feed source.
type Feed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Lang        string `json:"lang"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// ArticleFilters holds optional query parameters for filtering articles.
type ArticleFilters struct {
	Source            string
	Category          string
	Lang              string
	Offset            int
	Limit             int
	Refresh           bool
	ExcludeSources    []string
	ExcludeCategories []string
}

// ArticlesData is the payload returned by the GET /api/articles endpoint.
type ArticlesData struct {
	Articles []Article       `json:"articles"`
	Metadata ArticleMetadata `json:"metadata"`
}

// ArticleMetadata holds pagination info for the articles response.
type ArticleMetadata struct {
	Offset  int  `json:"offset"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

// ApiResponse standardises the JSON response format across all endpoints.
type ApiResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

// TrackEvent represents a user interaction event sent by the extension.
// Event must be either "share" or "bookmark".
type TrackEvent struct {
	ArticleID string `json:"articleId"`
	Event     string `json:"event"`
}
