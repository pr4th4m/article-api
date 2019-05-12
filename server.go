package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pratz/nine-article-api/endpoint"
	"github.com/pratz/nine-article-api/logger"
	"go.uber.org/zap"
)

func main() {
	host := flag.String("host", "localhost", "Server address")
	port := flag.String("port", "8080", "Server address")
	loglevel := zap.LevelFlag("loglevel", zap.InfoLevel, "Control log level")
	flag.Parse()

	// Configure env
	l := logger.New(*loglevel)
	env := &endpoint.ArticleEnv{Log: l}

	// Define routes
	var router = mux.NewRouter()
	router.HandleFunc(endpoint.VOne("articles"),
		env.Create).Methods(http.MethodPost)
	router.HandleFunc(endpoint.VOne("articles/{id}"),
		env.Get).Methods(http.MethodGet)
	router.HandleFunc(endpoint.VOne("tags/{tag}/{date}"),
		env.Search).Methods(http.MethodGet)

	address := fmt.Sprintf("%s:%s", *host, *port)
	l.Infof("Server running on %s", address)

	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length",
		"Accept-Encoding", "X-CSRF-Token", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	log.Fatal(http.ListenAndServe(address,
		handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router)))
}
