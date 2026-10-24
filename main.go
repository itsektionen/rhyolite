package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var RhyoliteDir string

func main() {
	address := os.Getenv("RHYOLITE_WEB_ADDR")
	if address == "" {
		log.Panicln("No RHYOLITE_WEB_ADDR specified")
	}

	RhyoliteDir = os.Getenv("RHYOLITE_ROOT_DIR")
	if RhyoliteDir == "" {
		log.Panicln("No RHYOLITE_ROOT_DIR specified")
	}

	watchForFileChanges()
	defer closeWatchForFileChanges()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Rhyolite Server"))
	})

	r.Post("/v1/{path}/{token}", handleExecuteRequest)

	log.Println("Starting on address:", address)
	http.ListenAndServe(address, r)
}
