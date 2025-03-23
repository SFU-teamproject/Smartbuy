package models

import "time"

type Review struct {
	ID           int       `json:"id"`
	SmartphoneID int       `json:"smartphone_id"`
	UserID       int       `json:"user_id"`
	Rating       int       `json:"rating"`
	Comment      *string   `json:"comment,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
