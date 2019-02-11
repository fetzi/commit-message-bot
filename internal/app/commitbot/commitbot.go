package commitbot

import (
	"io/ioutil"

	"github.com/cheshir/logrustash"
	log "github.com/sirupsen/logrus"

	"github.com/karriereat/commit-message-bot/internal/pkg/config"
	"github.com/karriereat/commit-message-bot/internal/pkg/handler"
	"github.com/karriereat/commit-message-bot/internal/pkg/server"
	"github.com/karriereat/commit-message-bot/internal/pkg/notifier"
	"github.com/karriereat/commit-message-bot/internal/pkg/refstore"
)

// Bot struct
type Bot struct {
	server *server.Server
}

// NewBot creates a new bot instance and adds the route handlers
func NewBot(config *config.Config) *Bot {
	srv := server.NewServer(config.Server.Port)

	commitLogger := buildCommitLogger(config)
	notifier := notifier.NewSlackNotifier(config)
	refStore, err := refstore.NewBoltRefStore(config.Database.Path)

	if err != nil {
		log.Error(err)
	}

	handler := handler.NewGitlabHookHandler(config, commitLogger, notifier, refStore)

	srv.AddRoute("/hooks/commit-message", handler)

	return &Bot{srv}
}

// Run starts the bot
func (b *Bot) Run() {
	b.server.Run()
}

// builds the desired commit logger instance
func buildCommitLogger(config *config.Config) *log.Logger {
	commitLogger := log.New()

	if config.CommitLog.Type == "logstash" {
		commitLogger.SetFormatter(&log.JSONFormatter{})

		hook, _ := logrustash.NewHookWithFields(
			"udp",
			config.CommitLog.Server,
			"commit-message-bot",
			log.Fields{
				"service": config.CommitLog.Servicename,
			},
		)

		commitLogger.Hooks.Add(hook)
		commitLogger.SetOutput(ioutil.Discard)
	}

	return commitLogger
}
