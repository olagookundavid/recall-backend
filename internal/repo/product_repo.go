package repo

import (
	"context"
	"errors"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductModel struct {
	DB *pgxpool.Pool
}

func (m ProductModel) Insert(product *domain.Product) (string, error) {
	query := `
		INSERT INTO tracked_product (
			user_id, name, company_name, store_name, country, category, phone, url, date_purchased
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	args := []any{
		product.UserId,
		product.Name,
		product.Company,
		product.Store,
		product.Country,
		product.Category,
		product.Phone,
		product.Url,
		product.Date,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id string
	err := m.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m ProductModel) GetProduct(user_id string) (*domain.Product, error) {
	query := ` SELECT id, user_id, name, company_name, store_name, country,  category, phone, url, date_purchased FROM tracked_product
	WHERE user_id = $1`
	var product domain.Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, user_id).Scan(
		&product.Id,
		&product.UserId,
		&product.Name,
		&product.Company,
		&product.Store,
		&product.Country,
		&product.Category,
		&product.Phone,
		&product.Url,
		&product.Date)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &product, nil
}

func (m ProductModel) GetAllProductsPaginatedWithNotification(limit, offset int) ([]*domain.ProductWithToken, error) {
	query := `
	SELECT 
    tp.id, tp.user_id, tp.name, tp.company_name, tp.store_name, 
    tp.country, tp.category, tp.phone, tp.url, tp.date_purchased,
    COALESCE(nt.token, '') AS token
	FROM tracked_product tp
	LEFT JOIN (
    SELECT DISTINCT ON (user_id) user_id, token
    FROM notification_tokens
    ORDER BY user_id, date_updated DESC
	) nt ON tp.user_id = nt.user_id
	ORDER BY tp.id
	LIMIT $1 OFFSET $2;

	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.ProductWithToken
	for rows.Next() {
		var product domain.ProductWithToken
		if err := rows.Scan(
			&product.Id,
			&product.UserId,
			&product.Name,
			&product.Company,
			&product.Store,
			&product.Country,
			&product.Category,
			&product.Phone,
			&product.Url,
			&product.Date,
			&product.Token,
		); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

func (m ProductModel) GetProducts(user_id string) ([]*domain.Product, error) {

	query := `SELECT id, user_id, name, company_name, store_name, country,  category, phone, url, date_purchased FROM tracked_product
	WHERE user_id = $1; `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []*domain.Product{}
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.Id,
			&product.UserId,
			&product.Name,
			&product.Company,
			&product.Store,
			&product.Country,
			&product.Category,
			&product.Phone,
			&product.Url,
			&product.Date)

		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (m ProductModel) DeleteProduct(id, userID string) error {
	query := ` DELETE FROM tracked_product WHERE id = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, id, userID)
	return err
}

func (m ProductModel) DeleteAllProductForUser(userID string) error {
	query := ` DELETE FROM tracked_product WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, userID)
	return err
}
