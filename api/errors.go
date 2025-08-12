package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type genericError struct {
	Err string `json:"error"`
}

func wrap(e error) error {
	return genericError{Err: e.Error()}
}

func (e genericError) Error() string {
	return e.Err
}

type errInvalidPageURL struct {
	Err string `json:"error"`
}

func newErrInvalidPageURL(field string) errInvalidPageURL {
	return errInvalidPageURL{Err: "invalid page url: " + field}
}

func (e errInvalidPageURL) Error() string {
	return e.Err
}

type errMissingFieldPrefix struct {
	Err string `json:"error"`
}

func newErrMissingFieldPrefix(field string) errMissingFieldPrefix {
	return errMissingFieldPrefix{Err: "missing request field: " + field}
}

func (e errMissingFieldPrefix) Error() string {
	return e.Err
}

type errMissingParamPrefix struct {
	Err string `json:"error"`
}

func newErrMissingParamPrefix(field string) errMissingParamPrefix {
	return errMissingParamPrefix{Err: "missing query param: " + field}
}

func (e errMissingParamPrefix) Error() string {
	return e.Err
}

type errMarshallResponse struct {
	Err string `json:"error"`
}

func newErrMarshallResponse() errMarshallResponse {
	return errMarshallResponse{Err: "unable to write response"}
}

func (e errMarshallResponse) Error() string {
	return e.Err
}

type errUnmarshallRequest struct {
	Err string `json:"error"`
}

func newErrUnmarshallRequest() errUnmarshallRequest {
	return errUnmarshallRequest{Err: "unable to read request body"}
}

func (e errUnmarshallRequest) Error() string {
	return e.Err
}

// writeError emulates what http.Error does but uses json instead of text to represent the data
// this also ensures that all error responses follow the same structure
func writeError(w http.ResponseWriter, error error) {
	w.Header().Set("Content-Type", "application/json")

	var errInvalidPageURL errInvalidPageURL
	var errMissingFieldPrefix errMissingFieldPrefix
	var errMissingParamPrefix errMissingParamPrefix
	var errMarshallResponse errMarshallResponse
	var errUnmarshallRequest errUnmarshallRequest

	switch {
	case errors.As(error, &errInvalidPageURL):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(error, &errMissingFieldPrefix):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(error, &errMissingParamPrefix):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(error, &errMarshallResponse):
		w.WriteHeader(http.StatusInternalServerError)
	case errors.As(error, &errUnmarshallRequest):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		error = wrap(error)
	}

	body, err := json.Marshal(error)
	if err != nil {
		log.Println(err)
	}

	_, _ = w.Write(body)
}
