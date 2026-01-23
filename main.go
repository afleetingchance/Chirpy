package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/afleetingchance/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s\n", err)
		os.Exit(1)
	}

	apiCfg := &apiConfig{
		db:        database.New(db),
		platform:  platform,
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.getChirpByIdHandler)
	mux.Handle("POST /api/chirps", apiCfg.middlewareUserAuth(http.HandlerFunc(apiCfg.createChirpHandler)))
	mux.Handle("DELETE /api/chirps/{chirpId}", apiCfg.middlewareUserAuth(http.HandlerFunc(apiCfg.deleteChirpHandler)))

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.Handle("PUT /api/users", apiCfg.middlewareUserAuth(http.HandlerFunc(apiCfg.updateUserHandler)))
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)

	// Webhooks
	mux.Handle("POST /api/polka/webhooks", apiCfg.middlewareAPIAuth(http.HandlerFunc(apiCfg.polkaWebhookHandler), "polkaKey"))

	server.ListenAndServe()
}
