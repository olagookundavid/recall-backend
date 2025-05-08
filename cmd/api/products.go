package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"recall-app/cmd/dto"
	"recall-app/internal/domain"
	"recall-app/internal/services"
)

// c.Param("") //for path param
func (app *Application) GetProductFromQR(c *gin.Context) {

	qrCode := c.Query("qr_code")
	if qrCode == "" {
		app.badResponse(c, "qr code cannot be empty")
		return
	}

	QrProducts, err := services.GetProductFromQrCode(c, qrCode)
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

	id, err := app.Handlers.Products.Insert(product)
	if err != nil {
		app.ServerErrorResponse(c, err)
		return
	}

	app.background(func() {
		checkProductRecallInput := CheckProductRecallInput{
			productName:       product.Name,
			productId:         id,
			category:          product.Category,
			date:              "",
			url:               product.Url,
			notificationToken: "",
		}
		_, err := app.CheckProductRecall(checkProductRecallInput)
		if err != nil {
			app.Logger.Error(fmt.Sprint("couldn't get product recall data : ", err.Error()), nil)
		}
	})

	rsp := dto.ProductResponse{
		Id:       id,
		UserId:   product.UserId,
		Name:     product.Name,
		Store:    product.Store,
		Company:  product.Company,
		Date:     product.Date,
		Country:  product.Country,
		Category: product.Category,
		Phone:    product.Phone,
		Url:      product.Url,
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully inserted product", "data": rsp})

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

// func (app *Application) GetProductFromFda(c *gin.Context) {

// 	product := c.Query("product")
// 	if product == "" {
// 		app.badResponse(c, "product cannot be empty")
// 		return
// 	}
// 	category := c.Query("category")
// 	if product == "" {
// 		app.badResponse(c, "category cannot be empty")
// 		return
// 	}
// 	date := c.Query("date")
// 	if product == "" {
// 		app.badResponse(c, "category cannot be empty")
// 		return
// 	}

// 	checkProductRecallInput := CheckProductRecallInput{
// 		productName:       product,
// 		productId:         id,
// 		category:          product.Category,
// 		date:              "",
// 		url:               "",
// 		notificationToken: "",
// 	}

// 	RecallsProducts, err := app.CheckProductRecall(product, "iddd", category, date, "")
// 	if err != nil {
// 		app.ServerErrorResponse(c, fmt.Errorf("couldn't get product recall data"))
// 		return
// 	}

// 	if RecallsProducts == nil {
// 		RecallsProducts = &[]FDARecall{}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": *RecallsProducts,
// 	})
// }

func (app *Application) CheckAllProductRecall() {
	const limit = 30
	offset := 0

	for {
		products, err := app.Handlers.Products.GetAllProductsPaginatedWithNotification(limit, offset)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("error fetching products: %v", err), nil)
			break
		}
		if len(products) == 0 {
			app.Logger.Info("Fetched all products.", nil)
			break // no more data
		}

		// Process each product

		for _, product := range products {
			checkProductRecallInput := CheckProductRecallInput{
				date:              time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
				url:               product.Url,
				notificationToken: product.Token,
				productName:       product.Name,
				productId:         product.Id,
				category:          product.Category,
			}

			app.background(func() {
				const maxRetries = 3
				const retryDelay = 2 * time.Second

				var err error
				for attempt := 1; attempt <= maxRetries; attempt++ {
					_, err = app.CheckProductRecall(checkProductRecallInput)
					if err == nil {
						break
					}
					app.Logger.Error(fmt.Sprintf("attempt %d: couldn't get product recall data: %s", attempt, err.Error()), nil)
					time.Sleep(retryDelay)
				}
				if err != nil {
					app.Logger.Error("final failure after retries: "+err.Error(), nil)
				}
			})
		}

		offset += limit
	}
}

//you have a possible recall
//notication struct

// recalls, err := CheckProductRecall("Creamy Peanut Butter 16 oz")
// if err != nil {
// 	log.Fatalf("Error: %v", err)
// }

// if len(recalls) > 0 {
// 	fmt.Println("⚠️ Product has been recalled:")
// 	for _, r := range recalls {
// 		fmt.Printf("- %s (Reason: %s)\n", r.ProductDescription, r.ReasonForRecall)
// 	}
// } else {
// 	fmt.Println("✅ Product not found in recent FDA recalls.")
// }
