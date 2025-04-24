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

	app := &api.Application{
		Wg:         sync.WaitGroup{},
		Config:     *cfg,
		Logger:     log,
		TokenMaker: tokenMaker,
		Mailer: mailer.New(
			cfg.Smtp.Host,
			cfg.Smtp.Port,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Sender,
		),
		Handlers: handlers.NewHandlers(pool),
	}
	cronjobs(log, app)

	err = server.Serve(app)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}
}

/*
for recall keeps-

CREATE TABLE user_notifications (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    recall_id INT REFERENCES recalls(id) ON DELETE CASCADE,
    notified_at TIMESTAMP DEFAULT NOW(),

    UNIQUE (user_id, recall_id)
);


CREATE TABLE recall_sync_state (
    product_type TEXT PRIMARY KEY,
    last_synced_date DATE NOT NULL
);

INSERT INTO recall_sync_state (product_type, last_synced_date)
VALUES ('food', '2024-01-01')
ON CONFLICT (product_type) DO NOTHING;

*/
