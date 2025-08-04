package api

import (
	"fmt"
	"net/http"
)

const (
	errMissingFieldPrefix = "missing request field: "
	errMissingParamPrefix = "missing query param: "
	errMarshallResponse   = "unable to write response"
	errUnmarshallRequest  = "unable to read request body"
)

// writeError emulates what http.Error does but uses json instead of text to represent the data
// this also ensures that all error responses follow the same structure
func writeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, error)
}
