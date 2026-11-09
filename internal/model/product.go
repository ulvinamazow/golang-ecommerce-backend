package model

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CreatedBy   int64     `json:"created_by"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
