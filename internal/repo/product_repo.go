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

func (m ProductModel) Insert(product *domain.Product) error {
	query := ` INSERT INTO tracked_product (user_id, name, company_name, store_name, country,  category, phone, url, date_purchased) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	args := []any{product.UserId, product.Name, product.Company, product.Store, product.Country, product.Category, product.Phone, product.Url, product.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m ProductModel) GetProducts(user_id string) (*domain.Product, error) {
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
