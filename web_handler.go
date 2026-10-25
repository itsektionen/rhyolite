package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type ExecuteRequest struct {
	rock        string
	commandName string
	token       string
	githubToken string
	githubActor string
}

var alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")
var githubTokenRegex = regexp.MustCompile("^gh[a-z]_[a-zA-Z0-9]+$")

func (r *ExecuteRequest) validateGithubToken() error {
	if !githubTokenRegex.Match([]byte(r.githubToken)) {
		return fmt.Errorf("Github token is invalid format")
	}
	return nil
}

func (r *ExecuteRequest) validateGithubActor() error {
	if !alphaNumericRegex.Match([]byte(r.githubActor)) {
		return fmt.Errorf("Github actor contains invalid characters")
	}
	return nil

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
		githubToken: r.Header.Get("X-Github-Token"),
		githubActor: r.Header.Get("X-Github-Actor"),
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
