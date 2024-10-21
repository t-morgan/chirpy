package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/t-morgan/chirpy/internal/auth"
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

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}
	tokenId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}

	err = validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error validating body: %s", err), err)
		return
	}

	cleanedChirp := cleanChirp(params.Body)

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: tokenId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	authorID := query.Get("author_id")
	sort := query.Get("sort")

	var dbChirps []database.Chirp

	if authorID != "" {
		authorID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to parse author id", err)
			return
		}

		dbChirps, err = cfg.dbQueries.GetAllChirpsByUserID(r.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Chirps not found for author id", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps for author id", err)
			return
		}
	}

	if authorID == "" {
		dbChirpsResult, err := cfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
			return
		}
		dbChirps = dbChirpsResult
	}

	var chirps []Chirp
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if sort == "desc" {
		slices.Reverse(chirps)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading chirpID", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})

}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading chirpID", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}
	tokenId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirp", err)
		return
	}

	if dbChirp.UserID != tokenId {
		respondWithError(w, http.StatusForbidden, "Forbidden", errors.New("not authorized"))
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
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
