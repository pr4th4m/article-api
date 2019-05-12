package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pratz/nine-article-api/content"
	"github.com/pratz/nine-article-api/logger"
)

// ArticleEnv article endpoint environment
type ArticleEnv struct {
	Log logger.Logger
}

// Create new article
func (e *ArticleEnv) Create(res http.ResponseWriter, req *http.Request) {

	// Parse input
	fields := content.Fields{}
	if err := ParseJSON(req, &fields); err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusBadRequest,
			JSONError{Error: err.Error()})
		return
	}

	article, err := content.NewArticle(context.Background(), e.Log)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	// Create new article
	if err := article.Save(fields); err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	msg := map[string]string{
		"success": fmt.Sprintf("Article %s saved successfully", fields.ID)}
	RenderJSON(res, http.StatusOK, msg)
}

// Get new article
func (e *ArticleEnv) Get(res http.ResponseWriter, req *http.Request) {

	ID := mux.Vars(req)["id"]
	article, err := content.NewArticle(context.Background(), e.Log)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	// Get new article
	fields, err := article.Get(ID)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	RenderJSON(res, http.StatusOK, fields)
}

// Search article by tag
func (e *ArticleEnv) Search(res http.ResponseWriter, req *http.Request) {

	tag := mux.Vars(req)["tag"]
	date := mux.Vars(req)["date"]

	article, err := content.NewArticle(context.Background(), e.Log)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	// Search article
	result, err := article.SearchByTag(tag, date)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	RenderJSON(res, http.StatusOK, result)
}
