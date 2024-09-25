package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
		CleanedBody string `json:"cleaned_body"`
	}
	cleanedChirp := cleanChirp(params.Body)
	respondWithJSON(w, 200, response{CleanedBody: cleanedChirp})
}

func validateChirp(chirp string) error {
	if len(chirp) > 140 {
		return errors.New("chirp is too long")
	}

	return nil
}

func cleanChirp(chirp string) string {
	replacement := "****"

	splitStr := strings.Split(chirp, " ")
	for i, word := range splitStr {
		lowerStr := strings.ToLower(word)
		if lowerStr == "kerfuffle" || lowerStr == "sharbert" || lowerStr == "fornax" {
			splitStr[i] = replacement
		}

	}
	chirp = strings.Join(splitStr, " ")

	return chirp
}