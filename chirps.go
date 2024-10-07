package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"fmt"
	"errors"

	"github.com/google/uuid"
	"github.com/t-morgan/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	cleanedChirp := cleanChirp(params.Body)

	createChirpParams := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: params.UserID,
	}
	var dbChirp database.Chirp
	dbChirp, err = cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, 201, chirp)
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