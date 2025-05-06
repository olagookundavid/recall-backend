package api

import (
	"net/http"
	"recall-app/cmd/dto"
	"recall-app/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *Application) CreateRecallHandler(c *gin.Context) {
	var req dto.RecallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.badResponse(c, err.Error())
		return
	}

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	date, err := time.Parse(layout, req.Date)
	if err != nil {
		app.badResponse(c, err.Error())
		return
	}

	recall := &domain.Recalls{
		Id:             req.Id,
		UserId:         tokenPayload.UserId,
		Date:           date,
		RecallId:       req.RecallId,
		FdaDescription: req.FdaDescription,
	}

	err = app.Handlers.Recalls.Insert(recall)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	// app.background(func() {})
	c.JSON(http.StatusOK, gin.H{"message": "Successfully inserted recall"})

}

func (app *Application) GetRecallHandler(c *gin.Context) {

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	recalls, err := app.Handlers.Recalls.GetRecalls(tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	rsp := dto.ConvertToRecallResponse(recalls)

	c.JSON(http.StatusOK, gin.H{"data": rsp})

}

func (app *Application) GetPotRecallHandler(c *gin.Context) {

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	recalls, err := app.Handlers.Products.GetProductWithPotRecalls(tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}
	if recalls == nil {
		recalls = []*domain.ProductWithPotRecall{}
	}

	c.JSON(http.StatusOK, gin.H{"data": recalls})

}

func (app *Application) DeleteRecallHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		app.badResponse(c, "qr code cannot be empty")
		return
	}

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	err = app.Handlers.Recalls.DeleteRecall(id, tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted recall"})

}
