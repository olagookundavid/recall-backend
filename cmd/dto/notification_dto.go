package dto

import (
	"recall-app/internal/domain"
	"time"
)

type NotificationRequest struct {
	Token string `json:"token"`
}

type NotificationResponse struct {
	Id    string    `json:"id"`
	Token string    `json:"token"`
	Date  time.Time `json:"date"`
}

func ConvertToNotificationResponse(notificationTokens *domain.NotificationTokens) NotificationResponse {
	return NotificationResponse{
		Id:    notificationTokens.Id,
		Token: notificationTokens.Token,
		Date:  notificationTokens.DateUpdated,
	}
}
