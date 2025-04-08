package domain

import "time"

type Token struct {
	Token  string    `json:"token"`
	UserID string    `json:"user_id"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"scope"`
}
