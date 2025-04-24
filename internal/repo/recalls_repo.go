package repo

import (
	"context"
	"errors"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecallsModel struct {
	DB *pgxpool.Pool
}

func (m RecallsModel) Insert(recall *domain.Recalls) error {
	query := ` INSERT INTO fda_recalls (id, user_id, date_purchased) 
				VALUES ($1, $2, $3)`
	args := []any{recall.Id, recall.UserId, recall.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m RecallsModel) GetRecalls(user_id string) (*domain.Recalls, error) {
	query := ` SELECT id, user_id, date_purchased FROM fda_recalls
	WHERE user_id = $1`
	var recall domain.Recalls
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, user_id).Scan(
		&recall.Id,
		&recall.UserId,
		&recall.Date)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &recall, nil
}

func (m RecallsModel) DeleteRecall(id, userID string) error {
	query := ` DELETE FROM fda_recalls WHERE id = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, id, userID)
	return err
}

func (m RecallsModel) DeleteAllRecallsForUser(userID string) error {
	query := ` DELETE FROM fda_recalls WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, userID)
	return err
}
