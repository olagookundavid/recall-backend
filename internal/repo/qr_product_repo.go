package repo

import (
	"fmt"
	httpclient "recall-app/internal/client"

	"recall-app/cmd/dto"

	"github.com/gin-gonic/gin"
)

func GetProductFromQrCode(c *gin.Context, id string) (*dto.QrProductResponse, error) {

	client := httpclient.NewClient("https://world.openfoodfacts.org/api")

	QrProduct := &dto.QrProductResponse{}

	path := fmt.Sprintf("/v0/product/%s.json", id)
	err := client.Do(c, "GET", path, nil, QrProduct, nil)
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintln(QrProduct.Status, QrProduct.StatusVerbose, QrProduct.Code))
	return QrProduct, nil
}
