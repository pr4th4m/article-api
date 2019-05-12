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

type contentArticle struct {
	esclient es.ElasticSearch
	log      logger.Logger
}

func (a *contentArticle) Save(fields Fields) error {
	_, err := a.esclient.Index(index, typ, fields.ID, fields)
	if err != nil {
		return err
	}
	a.log.Infof("Article %s saved successfully", fields.ID)
	return nil
}

func (a *contentArticle) SearchByTag(tag, date string) (TagSearchResult, error) {
	tagResult := TagSearchResult{}

	t, err := time.Parse(inputDatelayout, date)
	if err != nil {
		return tagResult, fmt.Errorf("Invalid date format")
	}
	formattedDate := t.Format(defaultDatelayout)
	a.log.Debugf("Date format changed to %s", formattedDate)

	searchResult, err := a.esclient.SearchByTag(index, typ, tag, formattedDate)
	if err != nil {
		return tagResult, err
	}

	aggs := searchResult.Aggregations
	idCount, _ := aggs.Terms(es.IDCountFieldName)
	relatedTags, _ := aggs.Terms(es.RelatedTagsFieldName)

	idSlice := []string{}
	for _, item := range idCount.Buckets {
		idSlice = append(idSlice, string(item.KeyNumber))
	}
	a.log.Debugf("Article id count %s", idSlice)

	relatedTagSlice := []string{}
	for _, item := range relatedTags.Buckets {
		relatedTagSlice = append(relatedTagSlice, item.Key.(string))
	}
	a.log.Debugf("Related tags list %s", relatedTagSlice)

	tagResult.Tag = tag
	tagResult.Count = searchResult.TotalHits()
	tagResult.Articles = idSlice
	tagResult.RelatedTags = relatedTagSlice

	a.log.Infof("Article search with tag:date %s:%s completed", tag, date)
	return tagResult, nil
}

func (a *contentArticle) Get(ID string) (Fields, error) {

	fields := Fields{}
	doc, err := a.esclient.Get(index, typ, ID)
	if err != nil {
		return fields, err
	}

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
