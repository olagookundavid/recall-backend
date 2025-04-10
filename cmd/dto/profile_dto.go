package dto

import (
	"recall-app/internal/domain"
)

type ProfileRequest struct {
	Name    *string `json:"name"`
	Country *string `json:"country"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Dob     *string `json:"date_of_birth"`
	Url     *string `json:"url"`
}

type ProfileResponse struct {
	Message string              `json:"message"`
	User    domain.UserResponse `json:"user"`
}

type PasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type PasswordResponse struct {
	Message string `json:"message"`
}
