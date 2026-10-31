package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// UpdateMaterial godoc
// @Summary      Update a material
// @Description  Updates specific fields of a material by its ID (PATCH)
// @Tags         materials
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Material ID"
// @Param        material body      models.Material false "Fields to update"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /materials/{id} [patch]
// @Security ApiKeyAuth
// PATCH
func (h *MaterialHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	materialID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(materialID)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var payload struct {
		Name         *string  `json:"name"`
		Unit         *string  `json:"unit"`
		CurrentStock *float64 `json:"current_stock"`
		MinStock     *float64 `json:"min_stock"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON in the request body")
		return
	}

	// Query construction
	query := "UPDATE materials SET "
	var args []interface{}

	if payload.Name != nil {
		query += "name = ?, "
		args = append(args, *payload.Name)
	}
	if payload.Unit != nil {
		query += "unit = ?, "
		args = append(args, *payload.Unit)
	}
	if payload.CurrentStock != nil {
		query += "current_stock = ?, "
		args = append(args, *payload.CurrentStock)
	}
	if payload.MinStock != nil {
		query += "min_stock = ?, "
		args = append(args, *payload.MinStock)
	}
	if len(args) == 0 {
		h.respondWithError(w, http.StatusBadRequest, "No fields provided for update")
		return
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = ?"

	args = append(args, id)

	result, err := h.DB.Exec(query, args...)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error during update")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to verify update status")
		return
	}

	if rowsAffected == 0 {
		h.respondWithError(w, http.StatusNotFound, "Material not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Material updated successfully",
	})
}
