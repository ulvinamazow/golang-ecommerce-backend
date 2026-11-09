package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"Name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"update_at"`
}
