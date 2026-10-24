package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	LoadConfig()

	watchForFileChanges()
	defer closeWatchForFileChanges()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Rhyolite Server"))
	})

	r.Post("/v1/{path}/{name}/{token}", handleExecuteRequest)

	log.Println("Starting on address:", config.WebAddress)
	http.ListenAndServe(config.WebAddress, r)
}
