package es

/*
NOTE: This package is meant to be abstraction layer over original package
Advantages:
- Easily swap underling original package (elastic)
- Easily mock ElasticSearch interface for testing
- Create custom search engine satisfying ElasticSearch interface
*/

import (
	"context"
	"fmt"
	"os"

	"github.com/olivere/elastic/v7"
)

const (
	RelatedTagsFieldName = "related"
	IDCountFieldName     = "ids"
	defaultHost          = "127.0.0.1"

	tagFieldName  = "tags.keyword"
	dateFieldName = "date"
	IDFieldName   = "id.keyword"
	topHitSize    = 10
)

// ElasticSearch top level interface
type ElasticSearch interface {
	Index(name, typ, ID string, body interface{}) (*elastic.IndexResponse, error)
	Get(name, typ, ID string) (*elastic.GetResult, error)
	SearchByTag(index, typ, tag, date string) (*elastic.SearchResult, error)
}

// NOTE: No need to pass our logger here
// as we can use original packages logger
// see NewES below
// customES elasticsearch wrapper
type customES struct {
	client *elastic.Client
	ctx    context.Context
}

// Index new document
func (e *customES) Index(name, typ, ID string,
	body interface{}) (*elastic.IndexResponse, error) {

	doc, err := e.client.Index().Index(name).Type(typ).Id(ID).BodyJson(body).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Get new document
func (e *customES) Get(name, typ, ID string) (*elastic.GetResult, error) {

	doc, err := e.client.Get().Index(name).Type(typ).Id(ID).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// SearchByTag search documents by tag
func (e *customES) SearchByTag(index, typ, tag,
	date string) (*elastic.SearchResult, error) {

	// Create match queries for date and tag
	dateQuery := elastic.NewMatchQuery(dateFieldName, date)
	tagQuery := elastic.NewMatchQuery(tagFieldName, tag)

	// Create bool query
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(dateQuery, tagQuery)

	// create id count and related aggregations
	relatedTagAgg := elastic.NewTermsAggregation().Field(tagFieldName).ExcludeValues(tag)
	idCountAgg := elastic.NewTermsAggregation().Field(IDFieldName).Size(topHitSize)

	searchResult, err := e.client.Search().
		Index(index).
		Query(boolQuery).
		Sort(dateFieldName, false).
		Size(0).
		Aggregation(RelatedTagsFieldName, relatedTagAgg).
		Aggregation(IDCountFieldName, idCountAgg).
		Pretty(true).
		Do(e.ctx)

	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

// NewES custom wrapper around elastic package
func NewES(ctx context.Context) (ElasticSearch, error) {

	// TODO: Figure out a better alternative
	// Hack to fix network issue inside docker
	host := os.Getenv("ES_HOST")
	if host == "" {
		host = defaultHost
	}
	addr := fmt.Sprintf("http://%s:9200", host)

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(addr),

		// NOTE: We can turn on "elastic" packages debugging with env var
		// https://github.com/olivere/elastic/wiki/Logging
		// elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		// elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]]", 0)),
	)

	if err != nil {
		return nil, err
	}

	return &customES{
		client: client,
		ctx:    ctx,
	}, nil
}
