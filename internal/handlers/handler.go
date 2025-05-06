package handlers

import (
	"recall-app/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Users        repo.UserModel
	Tokens       repo.TokenModel
	Notification repo.NotifcationTokensModel
	Products     repo.ProductModel
	Recalls      repo.RecallsModel
	PotRecalls   repo.PotRecallModel
	Transaction  repo.TransactionModel
}

func NewHandlers(db *pgxpool.Pool) Handlers {
	return Handlers{
		Products:     repo.ProductModel{DB: db},
		Recalls:      repo.RecallsModel{DB: db},
		Users:        repo.UserModel{DB: db},
		Tokens:       repo.TokenModel{DB: db},
		Transaction:  repo.TransactionModel{DB: db},
		PotRecalls:   repo.PotRecallModel{DB: db},
		Notification: repo.NotifcationTokensModel{DB: db},
	}
}
