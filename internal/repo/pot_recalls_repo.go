package repo

import (
	"context"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PotRecallModel struct {
	DB *pgxpool.Pool
}

func (m PotRecallModel) Insert(recall *domain.PotFDARecall) error {
	query := `
		INSERT INTO pot_fda_recalls (
			status, country, product_type, recalling_firm, address,
			voluntary_mandated, initial_firm_notification, distribution_pattern,
			recall_number, product_description, product_quantity, reason_for_recall,
			recall_initiation_date, termination_date, report_date, code_info, id
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16, $17
		)
		RETURNING id`
	args := []any{
		recall.Status,
		recall.Country,
		recall.ProductType,
		recall.RecallingFirm,
		recall.Address1,
		recall.VoluntaryMandated,
		recall.InitialFirmNotification,
		recall.DistributionPattern,
		recall.RecallNumber,
		recall.ProductDescription,
		recall.ProductQuantity,
		recall.ReasonForRecall,
		recall.RecallInitiationDate,
		recall.TerminationDate,
		recall.ReportDate,
		recall.CodeInfo,
		recall.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&recall.ID)
}

func (m PotRecallModel) GetAll() ([]*domain.PotFDARecall, error) {
	query := `
		SELECT 
			id, status, country, product_type, recalling_firm, address,
			voluntary_mandated, initial_firm_notification, distribution_pattern,
			recall_number, product_description, product_quantity, reason_for_recall,
			recall_initiation_date, termination_date, report_date, code_info,
			created_at, updated_at
		FROM pot_fda_recalls`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recalls []*domain.PotFDARecall

	for rows.Next() {
		var recall domain.PotFDARecall
		err := rows.Scan(
			&recall.ID,
			&recall.Status,
			&recall.Country,
			&recall.ProductType,
			&recall.RecallingFirm,
			&recall.Address1,
			&recall.VoluntaryMandated,
			&recall.InitialFirmNotification,
			&recall.DistributionPattern,
			&recall.RecallNumber,
			&recall.ProductDescription,
			&recall.ProductQuantity,
			&recall.ReasonForRecall,
			&recall.RecallInitiationDate,
			&recall.TerminationDate,
			&recall.ReportDate,
			&recall.CodeInfo,
			&recall.CreatedAt,
			&recall.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recalls = append(recalls, &recall)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return recalls, nil
}

func (m ProductModel) GetProductWithPotRecalls(user_id string) ([]*domain.ProductWithPotRecall, error) {
	query := `
	SELECT 
		tp.id, tp.user_id, tp.name, tp.company_name, tp.store_name, tp.country,  
		tp.category, tp.phone, tp.url, tp.date_purchased,
		pr.id, pr.status, pr.country, pr.product_type, pr.recalling_firm, pr.address,
		pr.voluntary_mandated, pr.initial_firm_notification, pr.distribution_pattern,
		pr.recall_number, pr.product_description, pr.product_quantity, pr.reason_for_recall,
		pr.recall_initiation_date, pr.termination_date, pr.report_date, pr.code_info,
		pr.created_at, pr.updated_at
	FROM tracked_product tp
	INNER JOIN pot_fda_recalls pr ON tp.id = pr.id
	WHERE tp.user_id = $1
	ORDER BY tp.id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productsMap := make(map[string]*domain.ProductWithPotRecall)

	for rows.Next() {
		var (
			recall      domain.PotFDARecall
			mainproduct domain.Product
		)

		err := rows.Scan(
			&mainproduct.Id, &mainproduct.UserId, &mainproduct.Name, &mainproduct.Company, &mainproduct.Store, &mainproduct.Country,
			&mainproduct.Category, &mainproduct.Phone, &mainproduct.Url, &mainproduct.Date,
			&recall.ID, &recall.Status, &recall.Country, &recall.ProductType, &recall.RecallingFirm, &recall.Address1,
			&recall.VoluntaryMandated, &recall.InitialFirmNotification, &recall.DistributionPattern,
			&recall.RecallNumber, &recall.ProductDescription, &recall.ProductQuantity, &recall.ReasonForRecall,
			&recall.RecallInitiationDate, &recall.TerminationDate, &recall.ReportDate, &recall.CodeInfo,
			&recall.CreatedAt, &recall.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		product, exists := productsMap[mainproduct.Id]
		if !exists {
			product = &domain.ProductWithPotRecall{
				Id:            mainproduct.Id,
				UserId:        mainproduct.UserId,
				Name:          mainproduct.Name,
				Company:       mainproduct.Company,
				Store:         mainproduct.Store,
				Country:       mainproduct.Country,
				Category:      mainproduct.Category,
				Phone:         mainproduct.Phone,
				Url:           mainproduct.Url,
				Date:          mainproduct.Date,
				PotFDARecalls: []*domain.PotFDARecall{},
			}
			productsMap[mainproduct.Id] = product
		}

		// Avoid appending empty recall if there was no match
		if recall.ID != "" {
			product.PotFDARecalls = append(product.PotFDARecalls, &recall)
		}
	}

	var products []*domain.ProductWithPotRecall
	for _, p := range productsMap {
		products = append(products, p)
	}
	return products, nil
}

func (m PotRecallModel) DeleteByID(id string) error {
	query := `DELETE FROM pot_fda_recalls WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, id)
	return err
}

func (m PotRecallModel) DeleteAllProductForUser(id, userID string) error {
	query := ` DELETE FROM pot_fda_recalls WHERE user_id = $1 AND `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, userID)
	return err
}
