package dto

type QrProductResponse struct {
	Code string `json:"code"`
	// also a pointer — so you can check if product is nil(4 jae)
	Product       *Product `json:"product"`
	StatusVerbose string   `json:"status_verbose"`
	Status        int      `json:"status"`
}

type Product struct {
	BrandOwner         string `json:"brand_owner"`
	BrandOwnerImported string `json:"brand_owner_imported"`
	Brands             string `json:"brands"`
	Countries          string `json:"countries"`
	ProductType        string `json:"product_type"`
	ProductName        string `json:"product_name"`
	ProductNameEn      string `json:"product_name_en"`
	ImageURL           string `json:"image_url"`
	FoodGroups         string `json:"food_groups"`
	CreatedT           int64  `json:"created_t"`
}
