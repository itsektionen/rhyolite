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

func HandleExecuteRequest(w http.ResponseWriter, r *http.Request) {
	path_var := chi.URLParam(r, "path")
	log.Println("Attempting to run: ", path_var)
	token := r.Header.Get("X-Deploy-Token")
	request := ExecuteRequest{
		directory: path_var,
		token:     token,
	}
	err := VerifyRequest(request)
	if err != nil {
		log.Println("Bad request: ", err.Error())
		w.WriteHeader(402)
		w.Write([]byte("Bad request, see the Rhyolite console for more info."))
		return
	}
	config, err := ReadConfig(request)
	if err != nil {
		log.Println("Unable to execute request:", err.Error())
		w.WriteHeader(403)
		w.Write([]byte("Unauthorized request, see the Rhyolite console for more info."))
		return
	}

	output, err := Execute(config, request)
	if err != nil {
		log.Println("Unable to execute:", err.Error())
		log.Println("Command history:\n", string(output))
		w.WriteHeader(500)
		w.Write([]byte("Server error, see the Rhyolite console for more info."))
		return
	}

	log.Println("Command executed successfully, history:\n", string(output))

	w.Write([]byte("deployed"))
}
