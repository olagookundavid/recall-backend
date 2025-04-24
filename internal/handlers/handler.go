package handlers

import (
	"recall-app/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Users    repo.UserModel
	Tokens   repo.TokenModel
	Fda      repo.SyncModel
	Products repo.ProductModel
	Recalls  repo.RecallsModel
}

func NewHandlers(db *pgxpool.Pool) Handlers {
	return Handlers{
		Products: repo.ProductModel{DB: db},
		Recalls:  repo.RecallsModel{DB: db},
		Users:    repo.UserModel{DB: db},
		Tokens:   repo.TokenModel{DB: db},
	}
}
