package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/r4stypasty/gobuild-api/internal/models"
)

type MaterialHandler struct {
	DB *sql.DB
}

func (h *MaterialHandler) respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

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

// Transactions
func (h *MaterialHandler) WriteOffMaterial(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	materialID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(materialID)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var payload struct {
		Quantity float64 `json:"quantity"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON in the request body")
		return
	}

	if payload.Quantity <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "quantity cannot be negative or zero")
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer tx.Rollback()

	query := `UPDATE materials SET current_stock = current_stock - ? WHERE id = ? AND current_stock >= ?`

	res, err := tx.ExecContext(r.Context(), query, payload.Quantity, id, payload.Quantity)
	if err != nil {
		log.Printf("[ERROR] Failed to execute update: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Material not found or insufficient stock")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Material successfully written off"}`))
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
