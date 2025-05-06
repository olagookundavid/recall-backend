package domain

import "time"

type NotificationTokens struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	Token       string    `json:"token"`
	DateUpdated time.Time `json:"date_updated"`
}
