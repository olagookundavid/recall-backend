package domain

import "time"

type Product struct {
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

type ProductWithToken struct {
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
	Token    string    `json:"token"`
}
type ProductWithPotRecall struct {
	Id            string          `json:"id"`
	UserId        string          `json:"user_id"`
	Name          string          `json:"name"`
	Store         string          `json:"store"`
	Company       string          `json:"company"`
	Date          time.Time       `json:"date"`
	Country       string          `json:"country"`
	Category      string          `json:"category"`
	Phone         string          `json:"phone"`
	Url           string          `json:"url"`
	PotFDARecalls []*PotFDARecall `json:"pot_recall"`
}
