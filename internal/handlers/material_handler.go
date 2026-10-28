package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/r4stypasty/gobuild-api/internal/models"
)

type MaterialHandler struct {
	DB *sql.DB
}

func (h *MaterialHandler) respondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// GET /materials
func (h *MaterialHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query("SELECT id, name, unit, current_stock, min_stock FROM materials")
	if err != nil {
		log.Printf("[ERROR] Database query failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	materials := make([]models.Material, 0)

	for rows.Next() {
		var m models.Material
		err := rows.Scan(&m.ID, &m.Name, &m.Unit, &m.CurrentStock, &m.MinStock)
		if err != nil {
			log.Printf("[ERROR] Scan row failed: %v", err)
			h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		materials = append(materials, m)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ERROR] Rows iteration failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(materials); err != nil {
		log.Printf("[ERROR] JSON encoding failed: %v", err)
	}
}

// POST /materials
func (h *MaterialHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var m models.Material

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON format or empty request body")
		return
	}

	if validationErrors := m.Validate(); validationErrors != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": validationErrors})
		return
	}

	query := `INSERT INTO materials (name, unit, current_stock, min_stock) VALUES (?, ?, ?, ?);`
	result, err := h.DB.Exec(query, m.Name, m.Unit, m.CurrentStock, m.MinStock)
	if err != nil {
		log.Printf("[ERROR] Insert material failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error occurred while saving material")
		return
	}

	id, _ := result.LastInsertId()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Material created successfully",
		"id":      id,
	})
}
