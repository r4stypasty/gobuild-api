package models

import "unicode/utf8"

type Material struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	CurrentStock float64 `json:"current_stock"`
	MinStock     float64 `json:"min_stock"`
}

func (m *Material) Validate() map[string]string {
	errors := make(map[string]string)

	if utf8.RuneCountInString(m.Name) < 3 {
		errors["name"] = "The material name is required and must contain at least 3 characters"
	}
	if m.Unit == "" {
		errors["unit"] = "Unit of measurement cannot be empty"
	}

	if m.CurrentStock < 0 {
		errors["current_stock"] = "Current balance cannot be negative"
	}

	if m.MinStock < 0 {
		errors["min_stock"] = "The minimum balance cannot be negative"
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
