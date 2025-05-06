package services

import (
	"fmt"
	"os"
	httpclient "recall-app/internal/client"

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
	err := client.Do(c, "GET", path, nil, QrProduct, nil)
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintln(QrProduct.Status, QrProduct.StatusVerbose, QrProduct.Code))
	return QrProduct, nil
}
