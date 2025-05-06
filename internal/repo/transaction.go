package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionModel struct {
	DB *pgxpool.Pool
}

func (m *TransactionModel) BeginTx(c context.Context) (pgx.Tx, error) {
	return m.DB.Begin(c)
}
