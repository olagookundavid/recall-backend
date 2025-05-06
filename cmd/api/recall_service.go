package api

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"

	httpclient "recall-app/internal/client"
	"recall-app/internal/domain"
)

type FDARecall struct {
	Status                  string `json:"status"`
	Country                 string `json:"country"`
	ProductType             string `json:"product_type"`
	RecallingFirm           string `json:"recalling_firm"`
	Address1                string `json:"address_1"`
	VoluntaryMandated       string `json:"voluntary_mandated"`
	InitialFirmNotification string `json:"initial_firm_notification"`
	DistributionPattern     string `json:"distribution_pattern"`
	RecallNumber            string `json:"recall_number"`
	ProductDescription      string `json:"product_description"`
	ProductQuantity         string `json:"product_quantity"`
	ReasonForRecall         string `json:"reason_for_recall"`
	RecallInitiationDate    string `json:"recall_initiation_date"`
	TerminationDate         string `json:"termination_date"`
	ReportDate              string `json:"report_date"`
	CodeInfo                string `json:"code_info"`
}

func FdaRecallToDomainObject(productId string, fdaRecall FDARecall) *domain.PotFDARecall {

	return &domain.PotFDARecall{
		ID:                      productId,
		Status:                  fdaRecall.Status,
		Country:                 fdaRecall.Country,
		ProductType:             fdaRecall.ProductType,
		RecallingFirm:           fdaRecall.RecallingFirm,
		Address1:                fdaRecall.Address1,
		VoluntaryMandated:       fdaRecall.VoluntaryMandated,
		InitialFirmNotification: fdaRecall.InitialFirmNotification,
		DistributionPattern:     fdaRecall.DistributionPattern,
		RecallNumber:            fdaRecall.RecallNumber,
		ProductDescription:      fdaRecall.ProductDescription,
		ProductQuantity:         fdaRecall.ProductQuantity,
		ReasonForRecall:         fdaRecall.ReasonForRecall,
		RecallInitiationDate:    fdaRecall.RecallInitiationDate,
		TerminationDate:         fdaRecall.TerminationDate,
		ReportDate:              fdaRecall.ReportDate,
		CodeInfo:                fdaRecall.CodeInfo,
	}
}

type FDAResponse struct {
	Results []FDARecall `json:"results"`
}

type CheckProductRecallInput struct {
	productName       string
	productId         string
	category          string
	date              string
	url               string
	notificationToken string
}

func (data CheckProductRecallInput) _printCheckProductRecallInput() {
	println(data.productName, " ", data.productId, " ", data.category, " ", data.date, " ", data.url, " ", data.notificationToken, " ")
}

func (app *Application) CheckProductRecall(input CheckProductRecallInput) (*[]FDARecall, error) {

	fdaApiKey := os.Getenv("FDA_API_KEY")
	status := "Ongoing"
	query := url.QueryEscape(fmt.Sprintf(`product_description:"%s" AND status:"%s" AND report_date:[%s TO *]`, input.productName, status, input.date))
	if input.date == "" {
		query = url.QueryEscape(fmt.Sprintf(`product_description:"%s" AND status:"%s"`, input.productName, status))
	}
	fdaFoodBaseUrl := os.Getenv("FDA_FOOD_BASE_URL")
	fdaDrugBaseUrl := os.Getenv("FDA_DRUG_BASE_URL")
	if fdaFoodBaseUrl == "" || fdaDrugBaseUrl == "" {
		return nil, fmt.Errorf("couldn't load fda base url to use to make call")
	}
	baseUrl := func() string {
		if input.category == "food" {
			return fdaFoodBaseUrl
		}
		return fdaDrugBaseUrl
	}

	apiURL := fmt.Sprintf("%s?search=%s&limit=30", baseUrl(), query)
	println(apiURL)
	client := httpclient.NewClient(apiURL)

	var fdaResp FDAResponse

	err := client.Do(context.Background(), "GET", "", nil, &fdaResp, map[string]string{
		"Authorization": fmt.Sprint("Basic ", fdaApiKey),
	})
	if err != nil {
		app.Logger.Error(err.Error(), nil)
		return nil, err
	}
	if len(fdaResp.Results) == 0 {
		return nil, nil
	}

	// go func() {
	// 	notification.SendNotification(app.MessagingClient, c, []string{""}, "", messaging.Notification{
	// 		Title:    "",
	// 		Body:     "",
	// 		ImageURL: "",
	// 	})
	// }()
	err = app.ProductRecallTransaction(input.productId, fdaResp.Results)
	if err != nil {
		return nil, err
	}

	return &fdaResp.Results, nil
}

func (app *Application) ProductRecallTransaction(productId string, recalls []FDARecall) error {

	c := context.Background()
	tx, err := app.Handlers.Transaction.BeginTx(c)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(c)
			return
		}
		err = tx.Commit(c)
		if err != nil {
			return
		}
	}()
	var wg sync.WaitGroup

	wg.Add(len(recalls))
	for _, recall := range recalls {
		go func() {
			defer wg.Done()
			err = app.Handlers.PotRecalls.Insert(FdaRecallToDomainObject(productId, recall))
			if err != nil {
				println(err.Error())
				println("Began")
			}
		}()
	}
	wg.Wait()
	return nil
}

// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// func RandomString(length int) string {
// 	b := make([]byte, length)
// 	for i := range b {
// 		n, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(charset))))
// 		if err != nil {
// 			return ""
// 		}
// 		b[i] = charset[n.Int64()]
// 	}
// 	return string(b)
// }
