package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	errorExistCode = 1
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	TOKEN, exists := os.LookupEnv("TOKEN")
	if !exists {
		logger.Error("TOKEN env varibale is not set")
		os.Exit(errorExistCode)
	}

	pref := tele.Settings{
		Token:  TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		msg := fmt.Errorf("error starting bot: %w", err).Error()
		logger.Error(msg)
		os.Exit(errorExistCode)
	}

	b.Handle("/hello", func(c tele.Context) error {
		senderName := c.Chat().FirstName
		msg := fmt.Sprintf("Hello, %s", senderName)
		return c.Send(msg)
	})

	logger.Info("starting...")
	b.Start()
}
