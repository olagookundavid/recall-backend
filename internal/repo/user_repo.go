package repo

import (
	"context"
	"errors"
	"recall-app/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrRecordAlreadyExist = errors.New("already exists")
	ErrDuplicateEmail     = errors.New("duplicate email")
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m UserModel) Insert(user *domain.User) error {
	query := ` INSERT INTO users (name, country, phone, email, password_hash) 
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING id, created_at, version`
	args := []any{user.Name, user.Country, user.Phone, user.Email, user.Password.Hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*domain.User, error) {
	query := ` SELECT id, created_at, name, country, phone, url, date_of_birth, email, password_hash, isAdmin, version FROM users 
	WHERE email = $1`
	var user domain.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Country,
		&user.Phone,
		&user.Url,
		&user.Dob,
		&user.Email,
		&user.Password.Hash,
		&user.IsAdmin,
		&user.Version)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *domain.User) error {
	query := ` UPDATE users SET name = $1, email = $2, password_hash = $3, country = $4, phone = $5, date_of_birth = $6, 
	url = $7, version = version + 1 WHERE id = $8 AND version = $9 RETURNING version;`
	args := []any{user.Name, user.Email, user.Password.Hash, user.Country, user.Phone, user.Dob, user.Url, user.ID, user.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
