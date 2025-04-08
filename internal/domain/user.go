package domain

import (
	// "errors"
	// "recall-app/internal/domain"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Phone     string    `json:"phone"`
	Url       string    `json:"url"`
	Dob       time.Time `json:"date_of_birth"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	Version   int       `json:"-"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(email string) {
	// v.Check(email != "", "email", "must be provided")
	// v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(password string) {
	// v.Check(password != "", "password", "must be provided")
	// v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	// v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func (user *User) NewUserResponse() UserResponse {
	return UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Url:     user.Url,
		Country: user.Country,
		Dob:     dateToString(user.Dob),
	}
}

func dateToString(date time.Time) string {

	formattedDate := date.Format("02-01-2006")

	if formattedDate == "01-01-0001" {
		return ""
	}
	return formattedDate
}

type UserResponse struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
	Url     string `json:"url"`
	Dob     string `json:"date_of_birth"`
	Email   string `json:"email"`
}
