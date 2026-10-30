package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/r4stypasty/gobuild-api/internal/models"
)

// GET /materials
func (h *MaterialHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.QueryContext(r.Context(), "SELECT id, name, unit, current_stock, min_stock FROM materials")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(materials); err != nil {
		log.Printf("[ERROR] JSON encoding failed: %v", err)
	}
}

// Alerting for low stock levels
func (h *MaterialHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.QueryContext(r.Context(), "SELECT id, name, unit, current_stock, min_stock FROM materials WHERE current_stock <= min_stock")
	if err != nil {
		log.Printf("[ERROR] Database query failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	alerts := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var current_stock, min_stock float64
		var name, unit string
		if err := rows.Scan(&id, &name, &unit, &current_stock, &min_stock); err != nil {
			log.Printf("[ERROR] Failed to scan row: %v", err)
			h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		alerts = append(alerts, map[string]interface{}{
			"id":            id,
			"name":          name,
			"unit":          unit,
			"current_stock": current_stock,
			"min_stock":     min_stock,
		})
	}
	if err = rows.Err(); err != nil {
		log.Printf("[ERROR] Rows iteration failed: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
