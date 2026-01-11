package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afleetingchance/Chirpy/internal/auth"
	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/afleetingchance/Chirpy/internal/types"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error hashing password: %s", err))
	}

	user, err := cfg.db.CreateUser(
		req.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	respondWithJSON(w, 201, types.ConvertUserForResponse(user))
}
