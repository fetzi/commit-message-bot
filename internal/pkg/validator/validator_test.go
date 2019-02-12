package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

type ValidatorTest struct {
	message string
	err     error
}

var structureValidatorTests = []ValidatorTest{
	ValidatorTest{
		"Foo bar",
		errors.New("Subject-Body structure not present"),
	},
	ValidatorTest{
		"\n\nMessage",
		errors.New("Subject line cannot be empty"),
	},
	ValidatorTest{
		"Subject\n\n",
		errors.New("Body cannot be empty"),
	},
	ValidatorTest{
		"Subject\n\nMessage",
		nil,
	},
}

var subjectValidatorTests = []ValidatorTest{
	ValidatorTest{
		"foo bar",
		errors.New("Subject-Body structure not present"),
	},
	ValidatorTest{
		"foo bar\n\nmessage",
		errors.New("First character of subject must be upper case"),
	},
	ValidatorTest{
		"This is a too long subject line with more than 50 characters\n\nmessage",
		errors.New("Subject line is limited to 50 characters"),
	},
	ValidatorTest{
		"This is a subject!\n\nmessage",
		errors.New("Subject should not end with a punctuation mark"),
	},
	ValidatorTest{
		"This is a subject?\n\nmessage",
		errors.New("Subject should not end with a punctuation mark"),
	},
	ValidatorTest{
		"This is a subject.\n\nmessage",
		errors.New("Subject should not end with a punctuation mark"),
	},
	ValidatorTest{"A correct subject\n\nmessage", nil},
}

var bodyValidatorTests = []ValidatorTest{
	ValidatorTest{
		"Commit message",
		errors.New("Subject-Body structure not present"),
	},
	ValidatorTest{
		"Subject line \n\nBody with a too long content exceeding 72 characters characters characters",
		errors.New("Body lines should be wrapped after 72 characters"),
	},
	ValidatorTest{
		"Subject line \n\nBody\nBody with a too long content exceeding 72 characters characters characters",
		errors.New("Body lines should be wrapped after 72 characters"),
	},
	ValidatorTest{
		"Subject\n\nBody text",
		nil,
	},
}

var ticketNumberValidatorTests = []ValidatorTest{
	ValidatorTest{
		"KSWAT-1234 Subject\n\nMessage",
		errors.New("Ticket number in subject line is not allowed"),
	},
	ValidatorTest{
		"Subject line\n\nbody text with KSWAT-1234 in it",
		errors.New("Ticket number should be prefixed with either \"Resolves\" or \"See\" in its own line"),
	},
	ValidatorTest{
		"Subject line\n\nMessage",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nResolves: KDEVOPS-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nResolves: KSWAT-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nResolves KSWAT-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nSee: KSWAT-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nSee KSWAT-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nResolves: KBUKA-1234\nSee: KSWAT-1234",
		nil,
	},
	ValidatorTest{
		"Subject line\n\nMessage\n\nResolves: KBUKA-1234, KBUKA-5678",
		nil,
	},
}

func TestStructureValidator(t *testing.T) {
	validator := StructureValidator{}

	for _, test := range structureValidatorTests {
		commit := makeCommitWithMessage(test.message)

		err := validator.Validate(commit)

		assert.Equal(t, test.err, err)
	}
}

func TestSubjectValidator(t *testing.T) {
	validator := SubjectValidator{}

	for _, test := range subjectValidatorTests {
		commit := makeCommitWithMessage(test.message)

		err := validator.Validate(commit)

		assert.Equal(t, test.err, err)
	}
}

func TestBodyValidator(t *testing.T) {
	validator := BodyValidator{}

	for _, test := range bodyValidatorTests {
		commit := makeCommitWithMessage(test.message)

		err := validator.Validate(commit)

		assert.Equal(t, test.err, err)
	}
}

func TestTicketNumberValidator(t *testing.T) {
	validator := TicketNumberValidator{}

	for _, test := range ticketNumberValidatorTests {
		commit := makeCommitWithMessage(test.message)

		err := validator.Validate(commit)

		assert.Equal(t, test.err, err)
	}
}

func makeCommitWithMessage(message string) gitlab.Commit {
	return gitlab.Commit{Message: message}
}
