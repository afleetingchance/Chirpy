package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/afleetingchance/Chirpy/internal/auth"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Error getting authorization token: %s", err))
			return
		}

		userUuid, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Unauthorized: %s", err))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userUuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
