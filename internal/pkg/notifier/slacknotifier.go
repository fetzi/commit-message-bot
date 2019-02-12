package notifier

import (
	"fmt"

	"github.com/karriereat/commit-message-bot/internal/pkg/config"
	"github.com/nlopes/slack"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

// SlackNotifier struct
type SlackNotifier struct {
	client       *slack.Client
	emoji        string
	fallbackUser string
}

// NewSlackNotifier creates a new notifier instance
func NewSlackNotifier(config *config.Config) *SlackNotifier {
	return &SlackNotifier{
		slack.New(config.Slack.Token),
		config.Slack.IconEmoji,
		config.Slack.FallbackUser,
	}
}

// Send sends a slack message with the commit message violation
func (n SlackNotifier) Send(project gitlab.Project, commit gitlab.Commit, message string) error {
	user, err := n.client.GetUserByEmail(commit.Author.Email)
	recipient := n.fallbackUser

	if err == nil {
		recipient = fmt.Sprintf("@%s", user.Name)
	}

	text := "A git commit message violates the message guidlines"
	attachment := slack.Attachment{
		Title:     message,
		TitleLink: commit.URL,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Project",
				Value: project.PathWithNamespace,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Author",
				Value: fmt.Sprintf("%s (%s)", commit.Author.Name, commit.Author.Email),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Commit ID",
				Value: commit.ID,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Message",
				Value: commit.Message,
				Short: false,
			},
		},
		Color: "warning",
	}

	params := slack.PostMessageParameters{
		Username:  "Better Git Bot",
		IconEmoji: n.emoji,
	}

	_, _, err = n.client.PostMessage(
		recipient,
		slack.MsgOptionText(text, false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionPostMessageParameters(params),
	)

	return err
}
