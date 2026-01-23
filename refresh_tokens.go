package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afleetingchance/Chirpy/internal/auth"
	"github.com/afleetingchance/Chirpy/internal/database"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting authorization token: %s", err))
		return
	}

	refreshToken, err := cfg.db.GetToken(req.Context(), tokenString)
	if err != nil {

		respondWithError(w, 500, fmt.Sprintf("Error retrieving refresh token: %s", err))
		return
	}

	zero := database.RefreshToken{}
	if refreshToken == zero ||
		(refreshToken.ExpiresAt.Valid && time.Now().After(refreshToken.ExpiresAt.Time)) ||
		refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Unauthorized")
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret)

	if err := cfg.db.RevokeToken(req.Context(), refreshToken.Token); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error revoking refresh token: %s", err))
		return
	}

	respondWithJSON(w, 200, response{Token: newToken})
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting authorization token: %s", err))
		return
	}

	refreshToken, err := cfg.db.GetToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error retrieving refresh token: %s", err))
		return
	}

	zero := database.RefreshToken{}
	if refreshToken == zero || (refreshToken.ExpiresAt.Valid && time.Now().After(refreshToken.ExpiresAt.Time)) {
		respondWithError(w, 401, "Unauthorized")
	}

	if err := cfg.db.RevokeToken(req.Context(), refreshToken.Token); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error revoking refresh token: %s", err))
		return
	}

	w.WriteHeader(204)
}
