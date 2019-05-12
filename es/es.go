package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

const (
	RelatedTagsFieldName = "related"
	IDCountFieldName     = "ids"

	tagFieldName  = "tags.keyword"
	dateFieldName = "date.keyword"
	IDFieldName   = "id.keyword"
	topHitSize    = 10
)

type ElasticSearch interface {
	Index(name, typ, ID string, body interface{}) (*elastic.IndexResponse, error)
	Get(name, typ, ID string) (*elastic.GetResult, error)
	SearchByTag(index, typ, tag, date string) (*elastic.SearchResult, error)
}

type customES struct {
	client *elastic.Client
	ctx    context.Context
}

func (e *customES) Index(name, typ, ID string,
	body interface{}) (*elastic.IndexResponse, error) {

	doc, err := e.client.Index().Index(name).Type(typ).Id(ID).BodyJson(body).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *customES) Get(name, typ, ID string) (*elastic.GetResult, error) {

	doc, err := e.client.Get().Index(name).Type(typ).Id(ID).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *customES) SearchByTag(index, typ, tag,
	date string) (*elastic.SearchResult, error) {

	dateQuery := elastic.NewMatchQuery(dateFieldName, date)
	tagQuery := elastic.NewMatchQuery(tagFieldName, tag)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(dateQuery, tagQuery)

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

func NewES(ctx context.Context) (ElasticSearch, error) {

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
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
