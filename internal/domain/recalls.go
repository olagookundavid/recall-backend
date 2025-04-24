package domain

import "time"

type Recalls struct {
	Id     string    `json:"id"`
	UserId string    `json:"user_id"`
	Date   time.Time `json:"date"`
}
