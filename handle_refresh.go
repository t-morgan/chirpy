package main

import (
	"net/http"
	"time"

	"github.com/t-morgan/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil ||
		refreshToken.RevokedAt.Valid ||
		refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Not authorized", err)
		return
	}

	expiresIn := time.Duration(DefaultExpiresInSeconds) * time.Second
	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to make JWT", err)
		return
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, tokenResponse{
		Token: jwt,
	})
}
