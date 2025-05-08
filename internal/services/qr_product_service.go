package services

import (
	"fmt"
	"os"
	httpclient "recall-app/internal/client"
	"time"

	"recall-app/cmd/dto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GetProductFromQrCode(c *gin.Context, id string) (*dto.QrProductResponse, error) {

	godotenv.Load()
	baseUrl := os.Getenv("OPEN_FOOD_FACTS_URL")
	if baseUrl == "" {
		return nil, fmt.Errorf("couldn't load open fact url to use to make call")
	}

	client := httpclient.NewClient(baseUrl)

	QrProduct := &dto.QrProductResponse{}

	path := fmt.Sprintf("/v0/product/%s.json", id)
	var err error

	// app.background(func() {
	const maxRetries = 3
	const retryDelay = 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = client.Do(c, "GET", path, nil, QrProduct, nil)
		if err == nil {
			break
		}
		time.Sleep(retryDelay)
	}
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintln(QrProduct.Status, QrProduct.StatusVerbose, QrProduct.Code))
	return QrProduct, nil
}
