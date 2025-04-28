package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// This interface basically allows us to send anything, (not typed)
type Envelope map[string]interface{}

// This functioin will make for better api responses to the caller
func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// Use this function to extract any id from the request
func ReadIDParam(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return "", errors.New("invalid id parameter")
	}
	return id, nil
}
