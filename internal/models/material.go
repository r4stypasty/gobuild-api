package models

type Material struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	CurrentStock float64 `json:"current_stock"`
	MinStock     float64 `json:"minimum_stock"`
}
