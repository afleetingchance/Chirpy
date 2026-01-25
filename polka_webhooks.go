package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Bad request: %s", err))
		return
	}

	switch params.Event {
	case "user.upgraded":
		cfg.upgradeUser(w, req, params.Data.UserId)
		return
	default:
		w.WriteHeader(204)
		return
	}
}

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, req *http.Request, userIdString string) {
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing user ID: %s", err))
		return
	}

	_, err = cfg.db.UpdateUserIsChirpyRed(
		req.Context(),
		database.UpdateUserIsChirpyRedParams{
			IsChirpyRed: true,
			ID:          userId,
		},
	)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error updating user: %s", err))
		return
	}

	w.WriteHeader(204)
}
