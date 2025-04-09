package dto

import (
	"recall-app/internal/domain"
)

type ProfileRequest struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
	Country  *string `json:"country"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	Dob      *string `json:"date_of_birth"`
	Url      *string `json:"url"`
}

type ProfileResponse struct {
	User domain.UserResponse `json:"user"`
}
