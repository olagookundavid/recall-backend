package repo

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type FDADate time.Time

// func (d *FDADate) UnmarshalJSON(b []byte) error {
// 	var s string
// 	if err := json.Unmarshal(b, &s); err != nil {
// 		return err
// 	}
// 	t, err := time.Parse("20060102", s)
// 	if err != nil {
// 		return err
// 	}
// 	*d = FDADate(t)
// 	return nil
// }

// func (d FDADate) Time() time.Time {
// 	return time.Time(d)
// }

// type Recall struct {
// 	RecallNumber             string  `json:"recall_number"`  // unique
// 	Classification           string  `json:"classification"` // Class I, II, III
// 	Status                   string  `json:"status"`         // Ongoing, Completed, etc.
// 	ProductDescription       string  `json:"product_description"`
// 	ReasonForRecall          string  `json:"reason_for_recall"`
// 	CodeInfo                 string  `json:"code_info"`
// 	RecallingFirm            string  `json:"recalling_firm"`
// 	City                     string  `json:"city"`
// 	State                    string  `json:"state"`
// 	Country                  string  `json:"country"`
// 	DistributionPattern      string  `json:"distribution_pattern"`
// 	ProductQuantity          string  `json:"product_quantity"`
// 	ReportDate               FDADate `json:"report_date"`                // YYYYMMDD
// 	RecallInitiationDate     FDADate `json:"recall_initiation_date"`     // YYYYMMDD
// 	CenterClassificationDate FDADate `json:"center_classification_date"` // YYYYMMDD
// 	ProductType              string  `json:"product_type"`               // food, drug, device
// 	EventID                  string  `json:"event_id"`
// }

// type FDAResponse struct {
// 	Results []Recall `json:"results"`
// }
// type SyncModel struct {
// 	DB *pgxpool.Pool
// }

// func (m SyncModel) FetchAndStoreRecallsSince(fromDate time.Time) (latestDate time.Time, err error) {
// 	const limit = 100
// 	for skip := 0; ; skip += limit {
// 		url := fmt.Sprintf("https://api.fda.gov/food/enforcement.json?search=report_date:[%s+TO+NOW]&limit=%d&skip=%d",
// 			fromDate.Format("20060102"), limit, skip)

// 		resp, err := http.Get(url)
// 		if err != nil {
// 			return latestDate, err
// 		}
// 		defer resp.Body.Close()

// 		var data FDAResponse
// 		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 			return latestDate, err
// 		}

// 		if len(data.Results) == 0 {
// 			break
// 		}

// 		for _, recall := range data.Results {
// 			err = m.UpsertRecall(recall)
// 			if err != nil {
// 				return latestDate, err
// 			}
// 			if recall.ReportDate.Time().After(latestDate) {
// 				latestDate = recall.ReportDate.Time()
// 			}
// 		}
// 	}
// 	return latestDate, nil
// }

// func (m SyncModel) UpsertRecall(r Recall) error {

// 	_, err := m.DB.Exec(context.Background(), `
// 		INSERT INTO recalls (
// 			recall_number, classification, status, product_description, reason_for_recall,
// 			code_info, recalling_firm, city, state, country, distribution_pattern, product_quantity,
// 			report_date, recall_initiation_date, center_classification_date, product_type, event_id
// 		) VALUES (
// 			$1, $2, $3, $4, $5,
// 			$6, $7, $8, $9, $10, $11, $12,
// 			$13, $14, $15, $16, $17
// 		)
// 		ON CONFLICT (recall_number) DO UPDATE SET
// 			classification = EXCLUDED.classification,
// 			status = EXCLUDED.status,
// 			product_description = EXCLUDED.product_description,
// 			reason_for_recall = EXCLUDED.reason_for_recall,
// 			code_info = EXCLUDED.code_info,
// 			recalling_firm = EXCLUDED.recalling_firm,
// 			city = EXCLUDED.city,
// 			state = EXCLUDED.state,
// 			country = EXCLUDED.country,
// 			distribution_pattern = EXCLUDED.distribution_pattern,
// 			product_quantity = EXCLUDED.product_quantity,
// 			report_date = EXCLUDED.report_date,
// 			recall_initiation_date = EXCLUDED.recall_initiation_date,
// 			center_classification_date = EXCLUDED.center_classification_date,
// 			product_type = EXCLUDED.product_type,
// 			event_id = EXCLUDED.event_id;
// 	`, r.RecallNumber,
// 		r.Classification,
// 		r.Status,
// 		r.ProductDescription,
// 		r.ReasonForRecall,
// 		r.CodeInfo,
// 		r.RecallingFirm,
// 		r.City,
// 		r.State,
// 		r.Country,
// 		r.DistributionPattern,
// 		r.ProductQuantity,
// 		r.ReportDate.Time(),
// 		r.RecallInitiationDate.Time(),
// 		r.CenterClassificationDate.Time(),
// 		r.ProductType,
// 		r.EventID,
// 	)

// 	return err
// }

// func (m SyncModel) GetLastSyncedDate(productType string) (time.Time, error) {
// 	var date time.Time
// 	err := m.DB.QueryRow(context.Background(), `SELECT last_synced_report_date FROM recall_sync_state WHERE product_type=$1`, productType).Scan(&date)
// 	if err == pgx.ErrNoRows {
// 		return time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC), nil
// 	}
// 	return date, err
// }

// func (m SyncModel) UpdateLastSyncedDate(productType string, newDate time.Time) error {
// 	_, err := m.DB.Exec(context.Background(), `
// 		INSERT INTO recall_sync_state (product_type, last_synced_date)
// 		VALUES ($1, $2)
// 		ON CONFLICT (product_type)
// 		DO UPDATE SET last_synced_date = EXCLUDED.last_synced_date
// 	`, productType, newDate)
// 	return err
// }

// // CREATE TABLE food_recalls (
// //     id SERIAL PRIMARY KEY,
// //     recall_number TEXT UNIQUE NOT NULL,
// //     classification TEXT, -- Class I, II, III
// //     status TEXT,         -- Ongoing, Completed, Terminated, Removed
// //     product_description TEXT,
// //     reason_for_recall TEXT,
// //     code_info TEXT,
// //     recalling_firm TEXT,
// //     city TEXT,
// //     state TEXT,
// //     country TEXT,
// //     distribution_pattern TEXT,
// //     product_quantity TEXT,
// //     report_date DATE,
// //     recall_initiation_date DATE,
// //     center_classification_date DATE,
// //     product_type TEXT,       -- e.g., Food, Drug, Device
// //     event_id TEXT,

// //     created_at TIMESTAMP DEFAULT NOW(),
// //     updated_at TIMESTAMP DEFAULT NOW()
// // );
