package dto

import (
	"recall-app/internal/domain"
	"time"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type RegisterUserResponse struct {
	AccessToken           string              `json:"access_token"`
	AccessTokenExpiresAt  time.Time           `json:"access_token_expires_at"`
	RefreshToken          string              `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time           `json:"refresh_token_expires_at"`
	User                  domain.UserResponse `json:"user"`
}
