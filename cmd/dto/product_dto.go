package dto

import (
	"recall-app/internal/domain"
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

func ConvertToProductResponse(products []*domain.Product) []*ProductResponse {
	respProduct := make([]*ProductResponse, 0, len(products))

	for _, row := range products {
		product := &ProductResponse{
			Id:       row.Id,
			UserId:   row.UserId,
			Name:     row.Name,
			Store:    row.Store,
			Company:  row.Company,
			Date:     row.Date,
			Country:  row.Country,
			Category: row.Category,
			Phone:    row.Phone,
			Url:      row.Url,
		}
		respProduct = append(respProduct, product)
	}
	return respProduct
}
