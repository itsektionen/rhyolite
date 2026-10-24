package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ExecuteRequest struct {
	rock        string
	commandName string
	token       string
}

func handleExecuteRequest(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	commandName := chi.URLParam(r, "name")
	token := chi.URLParam(r, "token")
	log.Println("Attempting to run: ", path)
	request := ExecuteRequest{
		rock:        path,
		commandName: commandName,
		token:       token,
	}
	config, err := findCommandConfig(request)
	if err != nil {
		log.Println("Unable to find a matching command:", err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden request: " + err.Error()))
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

	log.Println("Command executed successfully")

	w.Write(output)
}
