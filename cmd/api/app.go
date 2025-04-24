package api

import (
	"recall-app/internal/handlers"
	"recall-app/internal/logger"
	"recall-app/internal/mailer"
	"recall-app/internal/token"
	"sync"
)

type Application struct {
	Handlers   handlers.Handlers
	Config     Config
	Logger     *logger.Logger
	Wg         sync.WaitGroup
	Mailer     mailer.Mailer
	TokenMaker token.Maker
}

type Config struct {
	Port  int
	Env   string
	Token struct {
		TokenKey             string
		AccessTokenDuration  string
		RefreshTokenDuration string
	}
	Db struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}
