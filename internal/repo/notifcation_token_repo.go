package repo

import (
	"context"
	"errors"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CREATE TABLE IF NOT EXISTS notification_tokens (
//     id UUID NOT NULL REFERENCES tracked_product ON DELETE CASCADE,
//     user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
//     token text NOT NULL,
//     date_updated DATE NOT NULL);

type NotifcationTokensModel struct {
	DB *pgxpool.Pool
}

func (m NotifcationTokensModel) Upsert(notificationToken *domain.NotificationTokens) error {
	query := `
	INSERT INTO notification_tokens (token, user_id, date_updated)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id)
	DO UPDATE SET token = EXCLUDED.token, date_updated = EXCLUDED.date_updated;`

	args := []any{notificationToken.Token, notificationToken.UserId, notificationToken.DateUpdated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m NotifcationTokensModel) GetNotificationTokens(user_id string) (*domain.NotificationTokens, error) {
	query := ` SELECT id, user_id, token, date_updated FROM notification_tokens
	WHERE user_id = $1`
	var notificationToken domain.NotificationTokens
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, user_id).Scan(
		&notificationToken.Id,
		&notificationToken.UserId,
		&notificationToken.Token,
		&notificationToken.DateUpdated)

	if err != nil {
		println(err.Error())
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			{
				println("got here")
				return nil, ErrRecordNotFound
			}
		default:
			{

				println("got here 2")
				return nil, err
			}
		}
	}
	return &notificationToken, nil
}

func (m NotifcationTokensModel) DeleteNotificationToken(id, userID string) error {
	query := ` DELETE FROM notification_tokens WHERE id = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, id, userID)
	return err
}
