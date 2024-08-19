package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/timurzdev/social-rating-bot/internal/infrastructure/sqlite"
	"github.com/timurzdev/social-rating-bot/internal/repository/repos"
	tele "gopkg.in/telebot.v3"
)

const (
	errorExistCode = 1
	migrationsPath = "migrations"
)

type Config struct {
	Token          string
	DSN            string
	MigrationsPath string
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := getConfig()

	db, err := sqlite.Dial(cfg.DSN)
	if err != nil {
		return fmt.Errorf("error starting postgres: %w", err)
	}

	err = migrateDB(cfg.DSN, cfg.MigrationsPath)
	if err != nil {
		return fmt.Errorf("error during migrations: %w", err)
	}

	_ = repos.NewRatingRepository(db)

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("error starting bot: %w", err)
	}

	b.Handle("/hello", func(c tele.Context) error {
		senderName := c.Chat().FirstName
		msg := fmt.Sprintf("Hello, %s", senderName)
		return c.Send(msg)
	})

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		b.Stop()
		close(stopped)
	}()

	logger.Info("starting...")
	b.Start()

	<-stopped

	logger.Info("stopped")
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func getConfig() Config {
	cfg := Config{}
	TOKEN, exists := os.LookupEnv("TOKEN")
	if exists {
		cfg.Token = TOKEN
	}

	DSN, exists := os.LookupEnv("DSN")
	if exists {
		cfg.DSN = DSN
	}

	migrationsPath, exists := os.LookupEnv("MIGRATIONS_PATH")
	if exists {
		cfg.MigrationsPath = migrationsPath
	}
	return cfg
}

func migrateDB(dsn, path string) error {
	if path == "" {
		return errors.New("failed migration: empty path")
	}

	if dsn == "" {
		return errors.New("failed migration: empty dsn")
	}

	m, err := migrate.New(path, dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil || errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
