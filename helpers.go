package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resBody, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error marshalling JSON: %s", err))
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBody)
}
