package main

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/timurzdev/social-rating-bot/internal/config"
	"github.com/timurzdev/social-rating-bot/internal/infrastructure/postgres"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/timurzdev/social-rating-bot/internal/repository/repos"
	tele "gopkg.in/telebot.v3"
)

const (
	errorExistCode = 1
	migrationsPath = "migrations"
)

type Rofl struct {
	rofl int64
}

func (r Rofl) Recipient() string {
	return fmt.Sprintf("%d", r.rofl)
}

func run() error {
	rofl := Rofl{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	db, err := postgres.Dial(cfg.DSN()) // SDELAT POSTGRES
	if err != nil {
		return fmt.Errorf("error starting postgres: %w", err)
	}

	err = migrateDB(cfg.DSN(), cfg.MigrationsPath)
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

	b.Handle("/timurzhon", func(c tele.Context) error {
		userID := c.Sender().ID
		rofl.rofl = userID
		msg := fmt.Sprintf("Your userID: %d", userID)
		return c.Send(msg)
	})

	b.Handle("/roflan", func(c tele.Context) error {
		return c.ForwardTo(rofl, "zdarova")
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

func migrateDB(dsn, path string) error {
	log.Println(path, dsn)
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
