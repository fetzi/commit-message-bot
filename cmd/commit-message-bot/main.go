package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/karriereat/commit-message-bot/internal/app/commitbot"
	"github.com/karriereat/commit-message-bot/internal/pkg/config"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	err := config.Init()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	c := config.Load()

	bot := commitbot.NewBot(c)

	bot.Run()
}
