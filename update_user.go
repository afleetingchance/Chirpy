package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afleetingchance/Chirpy/internal/auth"
	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/afleetingchance/Chirpy/internal/types"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userID := req.Context().Value("user_id").(uuid.UUID)

	var params parameters
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error hashing password: %s", err))
	}

	user, err := cfg.db.UpdateUser(
		req.Context(),
		database.UpdateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		},
	)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error updating user: %s", err))
		return
	}

	respondWithJSON(w, 200, types.ConvertUserForResponse(user))
}
