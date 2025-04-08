package repo

import (
	"context"
	"errors"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenModel struct {
	DB *pgxpool.Pool
}

func (m TokenModel) Insert(token *domain.Token) error {
	query := ` INSERT INTO tokens (token, user_id, email, expiry, scope) 
				VALUES ($1, $2, $3, $4, $5)`
	args := []any{token.Token, token.UserID, token.Email, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TokenModel) Get(email, token_value string) (*domain.Token, error) {
	query := ` SELECT token, user_id, email, expiry, scope FROM tokens 
	WHERE email = $1 and token = $2`
	var token domain.Token
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, email, token_value).Scan(
		&token.Token,
		&token.UserID,
		&token.Email,
		&token.Expiry,
		&token.Scope)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &token, nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID string) error {
	query := ` DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, scope, userID)
	return err
}

func (m TokenModel) DeleteAllExpiredTokens() error {
	query := ` DELETE FROM tokens WHERE tokens.expiry < NOW()`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query)

	return err
}
