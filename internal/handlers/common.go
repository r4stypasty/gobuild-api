package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type MaterialHandler struct {
	DB *sql.DB
}

func (h *MaterialHandler) respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}