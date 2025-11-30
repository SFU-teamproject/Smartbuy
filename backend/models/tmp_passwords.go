package models

import "time"

type TmpPassword struct {
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TmpRequest struct {
	Email string `json:"email"`
}
