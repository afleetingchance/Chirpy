package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/afleetingchance/Chirpy/internal/types"
	"github.com/google/uuid"
)

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
			UserID: req.Context().Value("user_id").(uuid.UUID),
		})
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error creating chirp: %s", err))
		}

		respondWithJSON(w, 201, types.ConvertChirpForResponse(chirp))
	}
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	sortParam := req.URL.Query().Get("sort")
	sortDirection := "asc"
	if sortParam != "" {
		sortDirection = sortParam
	}
	sort := "created_at_" + sortDirection

	userIdString := req.URL.Query().Get("author_id")
	var userId uuid.UUID
	if userIdString != "" {
		var err error
		userId, err = uuid.Parse(userIdString)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error parsing user ID: %s", err))
			return
		}
	}

	rawChirps, err := cfg.db.GetChirps(
		req.Context(),
		database.GetChirpsParams{
			Sort:   sort,
			UserID: userId,
		},
	)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error retrieving chirps: %s", err))
		return
	}

	var chirps []types.Chirp
	for _, rawChirp := range rawChirps {
		chirps = append(chirps, types.ConvertChirpForResponse(rawChirp))
	}

	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirpByIdHandler(w http.ResponseWriter, req *http.Request) {
	chirpIdString := req.PathValue("chirpId")

	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing chirp ID: %s", err))
		return
	}

	rawChirp, err := cfg.db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error retrieving chirp: %s", err))
		return
	}

	respondWithJSON(w, 200, types.ConvertChirpForResponse(rawChirp))
}
