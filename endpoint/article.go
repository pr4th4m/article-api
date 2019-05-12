package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pratz/nine-article-api/content"
	"github.com/pratz/nine-article-api/logger"
)

type ArticleEnv struct {
	Log logger.Logger
}

func (e *ArticleEnv) Create(res http.ResponseWriter, req *http.Request) {

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

func (e *ArticleEnv) Get(res http.ResponseWriter, req *http.Request) {

	ID := mux.Vars(req)["id"]
	article, err := content.NewArticle(context.Background(), e.Log)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	fields, err := article.Get(ID)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	RenderJSON(res, http.StatusOK, fields)
}

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

	result, err := article.SearchByTag(tag, date)
	if err != nil {
		e.Log.Error(err)
		RenderJSON(res, http.StatusInternalServerError,
			JSONError{Error: err.Error()})
		return
	}

	RenderJSON(res, http.StatusOK, result)
}
