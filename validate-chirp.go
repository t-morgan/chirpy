package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, fmt.Sprintf("Error decoding parameters: %s", err), http.StatusBadRequest)
		return
	}

	err = validateChirp(params.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating body: %s", err), http.StatusBadRequest)
		return
	}

	type response struct {
		Valid bool `json:"valid"`
	}
	responseBody := response{Valid: true}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshalling response: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func validateChirp(chirp string) error {
	if len(chirp) > 140 {
		return errors.New("chirp is too long")
	}

	return nil
}
