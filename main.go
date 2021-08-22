package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/handlers"
	"github.com/adigunhammedolalekan/rest-unit-testing-sample/repository"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
)

func main() {
	database := repository.Database{}
	repo := repository.New(database)
	handler := handlers.New(repo)


	router := chi.NewRouter()
	router.Route("/api/posts", func(r chi.Router) {
		r.Get("/", handler.GetPosts)
		r.Post("/", handler.CreatePostHandler)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("API server up at %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("API server failed: ", err)
	}
}
