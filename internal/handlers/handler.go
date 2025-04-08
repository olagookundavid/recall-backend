package handlers

import (
	"recall-app/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Users  repo.UserModel
	Tokens repo.TokenModel
}

func NewHandlers(db *pgxpool.Pool) Handlers {
	return Handlers{

		Users:  repo.UserModel{DB: db},
		Tokens: repo.TokenModel{DB: db},
	}
}
