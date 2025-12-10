package handler

import (
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	e := envelope{
		"error": message,
	}
	err := writeJSON(w, status, e, nil)
	if err != nil {
		log.Println(err)
	}
}

func notfoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered an error and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
