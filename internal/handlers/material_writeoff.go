package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// WriteOffMaterial godoc
// @Summary      Write off material
// @Description  Decreases the material stock using a transaction to prevent negative balance
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "Material ID"
// @Param        quantity body      object  true  "Quantity to write off (JSON: {"quantity": 10})"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /materials/{id}/write-off [post]
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
