package site

import (
	"log"
	"net/http"
)

func httpError(w http.ResponseWriter, status int, err error) {
	log.Println("ERROR:", status, err.Error())
	renderError(w, status)
}

func internalServerError(w http.ResponseWriter, err error) {
	httpError(w, 500, err)
}

func notFoundError(w http.ResponseWriter, err error) {
	httpError(w, 404, err)
}
