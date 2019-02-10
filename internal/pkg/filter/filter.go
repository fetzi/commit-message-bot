package filter

import (
	"strings"

	"github.com/karriereat/commit-message-bot/internal/pkg/config"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

// Filter defines a commit message filter interface
type Filter interface {
	Filter(ev gitlab.PushEventPayload) bool
}

func GetFilters(conf *config.Config) []Filter {
	return []Filter{
		BodyStartsWithFilter{
			conf.FilteredStartings,
		},
		EmailFilter{
			conf.FilteredEmails,
		},
		GroupFilter{
			conf.FilteredGroups,
		},
	}
}

// BodyStartsWithFilter defines a filter for commit message beginnings
type BodyStartsWithFilter struct {
	startings []string
}

// Filter filters events where the commit message starts with one of the given startings
func (b BodyStartsWithFilter) Filter(ev gitlab.PushEventPayload) bool {
	for _, starting := range b.startings {
		for _, commit := range ev.Commits {
			if strings.HasPrefix(commit.Message, starting) {
				return true
			}
		}
	}

	return false
}

// EmailFilter defines a filter for commit author e-mail addresses
type EmailFilter struct {
	emails []string
}

// Filter filters events that ware initiated by a given user, identified by his/her e-mail address
func (ef EmailFilter) Filter(ev gitlab.PushEventPayload) bool {
	for _, email := range ef.emails {
		for _, commit := range ev.Commits {
			if email == commit.Author.Email {
				return true
			}
		}
	}

	return false
}

// GroupFilter defines a filter for gitlab groups
type GroupFilter struct {
	groups []string
}

// Filter filters all events on projects that are excluded
func (gf GroupFilter) Filter(ev gitlab.PushEventPayload) bool {
	for _, group := range gf.groups {
		if strings.HasPrefix(ev.Project.PathWithNamespace, group) {
			return true
		}
	}

	return false
}
