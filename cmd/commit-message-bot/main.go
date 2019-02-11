package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/karriereat/commit-message-bot/internal/app/commitbot"
	"github.com/karriereat/commit-message-bot/internal/pkg/config"
)

var cfgFile string

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	flag.StringVar(&cfgFile, "c", "./config.toml", "path to config file")
}

func main() {
	flag.Parse()

	c, err := config.Load(cfgFile)

	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	level, err := log.ParseLevel(c.Server.LogLevel)

	if err != nil {
		log.Error(fmt.Sprintf("Invalid log level: %s", c.Server.LogLevel))
	} else {
		log.SetLevel(level)
	}

	bot := commitbot.NewBot(c)

	bot.Run()
}
