package main

import (
	"context"
	"fmt"
	"recall-app/cmd/api"
	"recall-app/internal/handlers"
	"recall-app/internal/logger"
	"recall-app/internal/mailer"
	"recall-app/internal/server"
	"recall-app/internal/token"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {

	mode := loadModeEnv()
	if mode {
		gin.SetMode(gin.ReleaseMode)
	}
	//Check version and exit
	displayVersion("version")

	// Initialize logger
	log := logger.GetLogger(logger.Options{
		IsProduction: mode,
		AppName:      "Recall-king",
		Environment:  "dev",
		TraceID:      "recall-app-id",
	})
	defer log.Sync()

	dbUrl := loadDbUrl(log)
	tokenDeets := loadTokenDetails(log)
	smtpDeets := loadSmtpDetails(log)

	cfg := flagSetup(dbUrl, tokenDeets, smtpDeets)

	ctx := context.Background()
	pool, err := openDB(*cfg, ctx)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}
	defer pool.Close()
	log.Info("database connection pool established", nil)

	expvarSetup()

	tokenMaker, err := token.NewPasetoMaker(cfg.Token.TokenKey)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create token maker: %w", err).Error(), nil)
	}

	// fcmClient, err := FirebaseInit(ctx, log)
	// if err != nil {
	// 	log.Fatal(fmt.Errorf("cannot firebase client: %w", err).Error(), nil)
	// }
	// println("done")

	app := &api.Application{
		Wg:              sync.WaitGroup{},
		Config:          *cfg,
		Logger:          log,
		TokenMaker:      tokenMaker,
		MessagingClient: nil,
		Mailer: mailer.New(
			cfg.Smtp.Host,
			cfg.Smtp.Port,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Sender,
		),
		Handlers: handlers.NewHandlers(pool),
	}
	cronjobs(app)

	err = server.Serve(app)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}
}
