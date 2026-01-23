package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("user_id").(uuid.UUID)
	chirpIdString := req.PathValue("chirpId")

	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing chirp ID: %s", err))
		return
	}

	chirp, err := cfg.db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error retrieving chirp: %s", err))
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, "Forbidden")
		return
	}

	if err = cfg.db.DeleteChirp(req.Context(), chirp.ID); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error deleting chirp: %s", err))
		return
	}

	w.WriteHeader(204)
}
