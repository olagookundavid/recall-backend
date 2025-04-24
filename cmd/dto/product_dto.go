package dto

import (
	"time"
)

type ProductRequest struct {
	Name     string `json:"name"`
	Store    string `json:"store"`
	Company  string `json:"company"`
	Date     string `json:"date"`
	Country  string `json:"country"`
	Category string `json:"category"`
	Phone    string `json:"phone"`
	Url      string `json:"url"`
}

type ProductResponse struct {
	Id       string    `json:"id"`
	UserId   string    `json:"user_id"`
	Name     string    `json:"name"`
	Store    string    `json:"store"`
	Company  string    `json:"company"`
	Date     time.Time `json:"date"`
	Country  string    `json:"country"`
	Category string    `json:"category"`
	Phone    string    `json:"phone"`
	Url      string    `json:"url"`
}
