package handler

import (
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
	hook *gitlab.Webhook
	filters []filter.Filter
	validators []validator.Validator
	commitLogger *log.Logger
	notifier notifier.Notifier
	refStore refstore.RefStore
}

// NewGitlabHookHandler creates a new GitlabHookHandler instance
func NewGitlabHookHandler(conf *config.Config, commitLogger *log.Logger, notifier notifier.Notifier, refStore refstore.RefStore) *GitlabHookHandler {
	hook, _ := gitlab.New(gitlab.Options.Secret(conf.Gitlab.Token))

	return &GitlabHookHandler{
		hook,
		filter.GetFilters(conf),
		validator.GetValidators(),
		commitLogger,
		notifier,
		refStore,
	}
}

// ServeHTTP handles the gitlab hook request
func (handler GitlabHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := handler.hook.Parse(r, gitlab.PushEvents)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(err)
		return
	}

	switch payload.(type) {
	case gitlab.PushEventPayload:
		event := payload.(gitlab.PushEventPayload)

		for _, filter := range handler.filters {
			if filter.Filter(event) {
				log.Debug("Event filtered")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		for _, commit := range event.Commits {

			if handler.refStore.Exists(commit.ID) {
				log.Debug(fmt.Sprintf("Skipping commit %s because it was already validated", commit.ID))
				continue
			}

			hasError := false
			for _, validator := range handler.validators {
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

			handler.refStore.Put(commit.ID, commit.ID)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
