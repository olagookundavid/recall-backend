package api

import (
	"net/http"
	"recall-app/cmd/dto"
	"recall-app/internal/domain"
	"recall-app/internal/repo"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *Application) CreateNotificationHandler(c *gin.Context) {
	var req dto.NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.badResponse(c, err.Error())
		return
	}

	if req.Token == "" {
		app.badResponse(c, "Token cannot be empty")
		return
	}
	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	notificationToken := &domain.NotificationTokens{
		UserId:      tokenPayload.UserId,
		Token:       req.Token,
		DateUpdated: time.Now(),
	}

	err = app.Handlers.Notification.Upsert(notificationToken)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully set notification token for user"})

}

func (app *Application) GetNotificationHandler(c *gin.Context) {

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	notificationTokens, err := app.Handlers.Notification.GetNotificationTokens(tokenPayload.UserId)
	if err != nil {
		if err == repo.ErrRecordNotFound {
			app.invalidCredentialsResponse(c)
			return
		}
		app.ServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Successfully retrieved notification token for user",
		"notification_token": dto.ConvertToNotificationResponse(notificationTokens)})

}

func (app *Application) DeleteNotificationHandler(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		app.badResponse(c, "Notification ID cannot be empty")
		return
	}

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	err = app.Handlers.Notification.DeleteNotificationToken(id, tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted notification token"})

}
