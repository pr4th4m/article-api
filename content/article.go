package content

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pratz/nine-article-api/es"
	"github.com/pratz/nine-article-api/logger"
)

const (
	inputDatelayout   = "20060102"
	defaultDatelayout = "2006-01-02"
	index             = "content"
	typ               = "article"
)

// contentArticle - Article is type of Content
type contentArticle struct {
	esclient es.ElasticSearch
	log      logger.Logger
}

// Save new article
func (a *contentArticle) Save(fields Fields) error {
	_, err := a.esclient.Index(index, typ, fields.ID, fields)
	if err != nil {
		return err
	}
	a.log.Infof("Article %s saved successfully", fields.ID)
	return nil
}

// SearchByTag search articles
func (a *contentArticle) SearchByTag(tag, date string) (TagSearchResult, error) {
	tagResult := TagSearchResult{}

	// Convert input date format to match with saved date
	t, err := time.Parse(inputDatelayout, date)
	if err != nil {
		return tagResult, fmt.Errorf("Invalid date format")
	}
	formattedDate := t.Format(defaultDatelayout)
	a.log.Debugf("Date format changed to %s", formattedDate)

	// Search for articles
	searchResult, err := a.esclient.SearchByTag(index, typ, tag, formattedDate)
	if err != nil {
		return tagResult, err
	}

	// Get search aggregations
	aggs := searchResult.Aggregations
	idCount, idCountExists := aggs.Terms(es.IDCountFieldName)
	relatedTags, relatedTagsExists := aggs.Terms(es.RelatedTagsFieldName)

	idSlice := []string{}
	if idCountExists {
		for _, item := range idCount.Buckets {
			idSlice = append(idSlice, string(item.KeyNumber))
		}
	}
	a.log.Debugf("Article id count %s", idSlice)

	relatedTagSlice := []string{}
	if relatedTagsExists {
		for _, item := range relatedTags.Buckets {
			relatedTagSlice = append(relatedTagSlice, item.Key.(string))
		}
	}
	a.log.Debugf("Related tags list %s", relatedTagSlice)

	tagResult.Tag = tag
	tagResult.Count = searchResult.TotalHits()
	tagResult.Articles = idSlice
	tagResult.RelatedTags = relatedTagSlice

	a.log.Infof("Article search with tag:date %s:%s completed", tag, date)
	return tagResult, nil
}

// Get article by ID
func (a *contentArticle) Get(ID string) (Fields, error) {

	fields := Fields{}
	doc, err := a.esclient.Get(index, typ, ID)
	if err != nil {
		return fields, err
	}

	// Unmarshal raw result to go struct
	raw, err := doc.Source.MarshalJSON()
	if err != nil {
		return fields, err
	}
	if err := json.Unmarshal(raw, &fields); err != nil {
		return fields, err
	}
	a.log.Debugf("Article %s unmarshalled successfully", ID)

	a.log.Infof("Article %s found successfully", ID)
	return fields, nil
}

// NewArticle for article related operations
func NewArticle(ctx context.Context,
	l logger.Logger) (Content, error) {

	client, err := es.NewES(ctx)
	if err != nil {
		return nil, err
	}
	return &contentArticle{
		esclient: client,
		log:      l,
	}, nil
}
