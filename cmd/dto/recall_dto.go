package dto

import (
	"recall-app/internal/domain"
	"time"
)

type RecallRequest struct {
	Id             string `json:"id"`
	RecallId       string `json:"recall_id"`
	FdaDescription string `json:"fda_description"`
	Date           string `json:"date"`
}

type RecallResponse struct {
	Id             string    `json:"id"`
	UserId         string    `json:"user_id"`
	RecallId       string    `json:"recall_id"`
	FdaDescription string    `json:"fda_description"`
	Date           time.Time `json:"date"`
}

func ConvertToRecallResponse(recalls []*domain.Recalls) []*RecallResponse {
	respRecall := make([]*RecallResponse, 0, len(recalls))

	for _, row := range recalls {
		recall := &RecallResponse{
			Id:             row.Id,
			UserId:         row.UserId,
			Date:           row.Date,
			RecallId:       row.RecallId,
			FdaDescription: row.FdaDescription,
		}
		respRecall = append(respRecall, recall)
	}
	return respRecall
}
