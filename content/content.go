package content

// Content - top level interface
type Content interface {
	Save(fields Fields) error
	Get(ID string) (Fields, error)
	SearchByTag(tag, date string) (TagSearchResult, error)
}

// Fields for request input args
type Fields struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

// TagSearchResult for search response
type TagSearchResult struct {
	Tag         string   `json:"tag"`
	Count       int64    `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}
