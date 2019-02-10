package handler

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/karriereat/commit-message-bot/internal/pkg/config"
	"github.com/karriereat/commit-message-bot/internal/pkg/filter"
	"github.com/karriereat/commit-message-bot/internal/pkg/validator"
	"github.com/karriereat/commit-message-bot/internal/pkg/notifier"
	"github.com/karriereat/commit-message-bot/internal/pkg/refstore"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

// GitlabHookHandler struct
type GitlabHookHandler struct {
	conf *config.Config
	commitLogger *log.Logger
	notifier notifier.Notifier
	refStore refstore.RefStore
}

// NewGitlabHookHandler creates a new GitlabHookHandler instance
func NewGitlabHookHandler(conf *config.Config, commitLogger *log.Logger, notifier notifier.Notifier, refStore refstore.RefStore) *GitlabHookHandler {
	return &GitlabHookHandler{
		conf,
		commitLogger,
		notifier,
		refStore,
	}
}

// ServeHTTP handles the gitlab hook request
func (handler GitlabHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hook, _ := gitlab.New(gitlab.Options.Secret(handler.conf.GitlabToken))

	payload, err := hook.Parse(r, gitlab.PushEvents)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("%+v", err)
		return
	}

	switch payload.(type) {
	case gitlab.PushEventPayload:
		event := payload.(gitlab.PushEventPayload)

		filters := filter.GetFilters(handler.conf)

		for _, filter := range filters {
			if filter.Filter(event) {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		validators := validator.GetValidators()

		for _, commit := range event.Commits {

			if handler.refStore.Exists(commit.ID) {
				continue
			}

			hasError := false
			for _, validator := range validators {
				err := validator.Validate(commit)

				if err != nil {
					handler.commitLogger.WithFields(log.Fields{
						"ref":     commit.ID,
						"project": event.Project.PathWithNamespace,
						"valid":   false,
					}).Info(err.Error())
					hasError = true

					notifierError := handler.notifier.Send(event.Project, commit, err.Error())

					if (notifierError != nil) {
						log.Error(notifierError)
					}
					break
				}
			}

			if !hasError {
				handler.commitLogger.WithFields(log.Fields{
					"requestId": commit.ID,
					"project":   event.Project.PathWithNamespace,
					"valid":     true,
				}).Info("Successfully validated commit")
			}

			handler.refStore.Put(commit.ID, event.Project.PathWithNamespace)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func respondSuccess(message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := make(map[string]string)
	data["message"] = message

	json, _ := json.Marshal(data)
	w.Write(json)
}
