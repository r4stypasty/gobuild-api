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

func (h *MaterialHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, unit, current_stock, min_stock FROM materials")
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	materials := make([]models.Material, 0)

	for rows.Next() {
		var m models.Material
		err := rows.Scan(&m.ID, &m.Name, &m.Unit, &m.CurrentStock, &m.MinStock)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		materials = append(materials, m)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(materials); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}
