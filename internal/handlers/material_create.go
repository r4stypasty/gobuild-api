package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/r4stypasty/gobuild-api/internal/models"
)

// POST /materials
func (h *MaterialHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var m models.Material

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON format or empty request body")
		return
	}

	if validationErrors := m.Validate(); validationErrors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": validationErrors})
		return
	}

	query := `INSERT INTO materials (name, unit, current_stock, min_stock) VALUES (?, ?, ?, ?);`
	result, err := h.DB.ExecContext(r.Context(), query, m.Name, m.Unit, m.CurrentStock, m.MinStock)
	if err != nil {
		log.Printf("[ERROR] Insert material failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error occurred while saving material")
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Material created successfully",
		"id":      id,
	})
}
