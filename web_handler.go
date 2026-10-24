package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ExecuteRequest struct {
	directory string
	token     string
}

func handleExecuteRequest(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	token := chi.URLParam(r, "token")
	log.Println("Attempting to run: ", path)
	request := ExecuteRequest{
		directory: path,
		token:     token,
	}
	config, err := findCommandConfig(request)
	if err != nil {
		log.Println("Unable to find a matching command:", err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden request, see the Rhyolite console for more info."))
		return
	}

	output, err := execute(config, request)
	if err != nil {
		log.Println("Unable to execute:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Command failed:\n"))
		w.Write(output)
		return
	}

	log.Printf("Command executed successfully, history:\n%s\n", string(output))

	w.Write(output)
}
