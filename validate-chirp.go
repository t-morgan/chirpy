package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	err = validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error validating body: %s", err), err)
		return
	}

	type response struct {
		Valid bool `json:"valid"`
	}
	respondWithJSON(w, 200, response{Valid: true})
}

func validateChirp(chirp string) error {
	if len(chirp) > 140 {
		return errors.New("chirp is too long")
	}

	return nil
}
