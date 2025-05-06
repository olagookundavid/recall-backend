package domain

import "time"

type PotFDARecall struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	Country                 string    `json:"country"`
	ProductType             string    `json:"product_type"`
	RecallingFirm           string    `json:"recalling_firm"`
	Address1                string    `json:"address_1"`
	VoluntaryMandated       string    `json:"voluntary_mandated"`
	InitialFirmNotification string    `json:"initial_firm_notification"`
	DistributionPattern     string    `json:"distribution_pattern"`
	RecallNumber            string    `json:"recall_number"`
	ProductDescription      string    `json:"product_description"`
	ProductQuantity         string    `json:"product_quantity"`
	ReasonForRecall         string    `json:"reason_for_recall"`
	RecallInitiationDate    string    `json:"recall_initiation_date"`
	TerminationDate         string    `json:"termination_date"`
	ReportDate              string    `json:"report_date"`
	CodeInfo                string    `json:"code_info"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
