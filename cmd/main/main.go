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

	var (
		ginMode = gin.DebugMode
	)
	if mode {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
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

	cfg := flagSetup(dbUrl, tokenDeets)

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

	app := &api.Application{
		Wg:         sync.WaitGroup{},
		Config:     *cfg,
		Logger:     log,
		TokenMaker: tokenMaker,
		Mailer: mailer.New(
			"smtp.gmail.com",
			587,
			"erijesudo@gmail.com",
			"whpemugxjgmincph",
			"erijesudo@gmail.com",
		),
		Handlers: handlers.NewHandlers(pool),
	}

	err = server.Serve(app)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}
	print("nil error")
}
