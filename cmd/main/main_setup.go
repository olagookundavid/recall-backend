package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"recall-app/internal/logger"
	"recall-app/internal/vcs"

	"recall-app/cmd/api"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"google.golang.org/api/option"
)

var (
	version = vcs.Version()
)

func expvarSetup() {
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}

func openDB(cfg api.Config, ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func displayVersion(flagStr string) {
	displayVersion := flag.Bool(flagStr, false, "Display version and exit")
	flag.Parse()
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}
}

func loadDbUrl(log *logger.Logger) string {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL env variable missing", nil)
	}
	return dbUrl
}

func loadModeEnv() bool {
	godotenv.Load()
	return (os.Getenv("IS_PROD")) == "true"
}

func loadTokenDetails(log *logger.Logger) map[string]string {

	godotenv.Load()
	tokenKey := os.Getenv("TOKEN_KEY")
	accessTokenDuration := os.Getenv("ACCESS_TOKEN_DURATION")
	refreshTokenDuration := os.Getenv("REFRESH_TOKEN_DURATION")

	if tokenKey == "" || accessTokenDuration == "" || refreshTokenDuration == "" {
		log.Fatal("Couldn't load token details", nil)
	}
	tokenMap := map[string]string{
		"token_key":              tokenKey,
		"access_token_duration":  accessTokenDuration,
		"refresh_token_duration": refreshTokenDuration,
	}
	return tokenMap
}

func loadSmtpDetails(log *logger.Logger) map[string]string {

	godotenv.Load()
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpSender := os.Getenv("SMTP_SENDER")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || smtpSender == "" {
		log.Fatal("Couldn't load smtp details", nil)
	}
	smtpMap := map[string]string{
		"smtp_host":     smtpHost,
		"smtp_port":     smtpPort,
		"smtp_username": smtpUsername,
		"smtp_password": smtpPassword,
		"smtp_sender":   smtpSender,
	}
	return smtpMap
}

func flagSetup(dbUrl string, tokenDeets map[string]string, smtpDeets map[string]string) *api.Config {

	var cfg api.Config

	//env and port
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	//db and settings
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", dbUrl, "PostgreSQL DSN")
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	//tokenDeets
	flag.StringVar(&cfg.Token.TokenKey, "token-key", tokenDeets["token_key"], "Token Key")
	flag.StringVar(&cfg.Token.AccessTokenDuration, "access-token-duration", tokenDeets["access_token_duration"], "Access Token Duration")
	flag.StringVar(&cfg.Token.RefreshTokenDuration, "refresh-token-duration", tokenDeets["refresh_token_duration"], "Refresh Token Duration")

	//smpt
	port, _ := strconv.Atoi(smtpDeets["smtp_port"])
	flag.StringVar(&cfg.Smtp.Host, "smtp-host", smtpDeets["smtp_host"], "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", port, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", smtpDeets["smtp_username"], "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", smtpDeets["smtp_password"], "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", smtpDeets["smtp_sender"], "SMTP sender")
	return &cfg
}

func cronjobs(app *api.Application) {
	c := cron.New()

	// Run every sat at 12noon
	c.AddFunc("0 12 * * 6", func() {
		app.CheckAllProductRecall()
	})
	app.Logger.Info("Starting scheduler...", nil)
	c.Start()
}

func FirebaseInit(ctx context.Context, log *logger.Logger) (*messaging.Client, error) {
	path, err := filepath.Abs("recall-king-firebase-adminsdk-key.json")
	if err != nil {
		log.Fatal("Cannot get firebase json path", nil)
	}
	// Use the path to your service account credential json file
	opt := option.WithCredentialsFile(path)
	// Create a new firebase app
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	// Get the FCM object
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	println("Client initailized")
	return fcmClient, nil
}
