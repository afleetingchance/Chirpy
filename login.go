package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afleetingchance/Chirpy/internal/auth"
	"github.com/afleetingchance/Chirpy/internal/types"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	canLogin, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error verifying password: %s", err))
		return
	}

	if !canLogin {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	respondWithJSON(w, 200, types.ConvertUserForResponse(user))
}
