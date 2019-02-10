package notifier

import "gopkg.in/go-playground/webhooks.v5/gitlab"

// Notifier interface
type Notifier interface {
	Send(project gitlab.Project, commit gitlab.Commit, message string) error
}
