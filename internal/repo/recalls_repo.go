package repo

import (
	"context"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecallsModel struct {
	DB *pgxpool.Pool
}

func (m RecallsModel) Insert(recall *domain.Recalls) error {
	query := ` INSERT INTO fda_recalls (id, user_id, fda_description, recall_id, date_recalled) 
				VALUES ($1, $2, $3, $4, $5)`
	args := []any{recall.Id, recall.UserId, recall.FdaDescription, recall.RecallId, recall.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m RecallsModel) GetRecalls(user_id string) ([]*domain.Recalls, error) {
	query := ` SELECT id, user_id, fda_description, recall_id, date_recalled FROM fda_recalls
	WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	recalls := []*domain.Recalls{}
	for rows.Next() {
		var recall domain.Recalls
		err := rows.Scan(
			&recall.Id,
			&recall.UserId,
			&recall.FdaDescription,
			&recall.RecallId,
			&recall.Date)

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
