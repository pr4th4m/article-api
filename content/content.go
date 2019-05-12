package content

type Content interface {
	Save(fields Fields) error
	Get(ID string) (Fields, error)
	SearchByTag(tag, date string) (TagSearchResult, error)
}

type Fields struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type TagSearchResult struct {
	Tag         string   `json:"tag"`
	Count       int64    `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}
