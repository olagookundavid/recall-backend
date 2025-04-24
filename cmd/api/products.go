package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"recall-app/cmd/dto"
	"recall-app/internal/domain"
	"recall-app/internal/repo"
)

// c.Param("") //for path param
func (app *Application) GetProductFromQR(c *gin.Context) {

	qrCode := c.Query("qr_code")
	if qrCode == "" {
		app.badResponse(c, "qr code cannot be empty")
		return
	}

	QrProducts, err := repo.GetProductFromQrCode(c, qrCode)
	if err != nil {
		app.ServerErrorResponse(c, fmt.Errorf("couldn't get qr code data"))
		return
	}

	if QrProducts.Status == 0 || QrProducts.Product == nil {
		app.badResponse(c, "product not found from qr code")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": QrProducts,
	})
}

func (app *Application) CreateProductHandler(c *gin.Context) {
	var req dto.ProductRequest
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

	product := &domain.Product{
		Name:     req.Name,
		Phone:    req.Phone,
		Country:  req.Country,
		UserId:   tokenPayload.UserId,
		Store:    req.Store,
		Company:  req.Company,
		Date:     date,
		Category: req.Category,
		Url:      req.Url,
	}

	err = app.Handlers.Products.Insert(product)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	// rsp := dto.ProductResponse{
	// 	Id:       product.Id,
	// 	UserId:   product.UserId,
	// 	Name:     product.Name,
	// 	Store:    product.Store,
	// 	Company:  product.Company,
	// 	Date:     product.Date,
	// 	Country:  product.Country,
	// 	Category: product.Category,
	// 	Phone:    product.Phone,
	// 	Url:      product.Url,
	// }
	c.JSON(http.StatusOK, gin.H{"message": "Successfully inserted product"})

}

func (app *Application) GetProductHandler(c *gin.Context) {

	tokenPayload, err := getTokenPayloadFromContext(c)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	products, err := app.Handlers.Products.GetProducts(tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	rsp := dto.ConvertToProductResponse(products)

	c.JSON(http.StatusOK, gin.H{"data": rsp})

}

func (app *Application) DeleteProductHandler(c *gin.Context) {
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

	err = app.Handlers.Products.DeleteProduct(id, tokenPayload.UserId)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted product"})

}
