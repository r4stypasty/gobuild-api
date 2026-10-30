package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DeleteMaterial godoc
// @Summary      Delete a material
// @Description  Deletes a material by its ID
// @Tags         materials
// @Param        id   path      int  true  "Material ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /materials/{id} [delete]
func (h *MaterialHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	materialID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(materialID)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	query := `DELETE FROM materials WHERE id = ?`
	res, err := h.DB.ExecContext(r.Context(), query, id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		h.respondWithError(w, http.StatusNotFound, "Material not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
