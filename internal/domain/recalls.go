package domain

import "time"

type Recalls struct {
	Id             string    `json:"id"`
	UserId         string    `json:"user_id"`
	RecallId       string    `json:"recall_id"`
	FdaDescription string    `json:"fda_description"`
	Date           time.Time `json:"date"`
}
