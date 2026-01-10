package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/google/uuid"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	var params parameters
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		words := strings.Split(params.Body, " ")
		for i, word := range words {
			if _, ok := badWords[strings.ToLower(word)]; ok {
				words[i] = "****"
			}
		}
		cleaned := strings.Join(words, " ")

		chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: params.UserId,
		})
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error creating chirp: %s", err))
		}

		respondWithJSON(w, 201, convertChirpForResponse(chirp))
	}
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	sort := "created_at ASC"

	rawChirps, err := cfg.db.GetChirps(req.Context(), sort)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error retrieving chirps: %s", err))
		return
	}

	var chirps []Chirp
	for _, rawChirp := range rawChirps {
		chirps = append(chirps, convertChirpForResponse(rawChirp))
	}

	respondWithJSON(w, 200, chirps)
}

func convertChirpForResponse(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}
}
