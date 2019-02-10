package validator

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

// Validator defines a commit message validator interface
type Validator interface {
	Validate(c gitlab.Commit) error
}

// GetValidators gets a slice of all available validators
func GetValidators() []Validator {
	return []Validator{
		StructureValidator{},
		SubjectValidator{},
		BodyValidator{},
		TicketNumberValidator{},
	}
}

// StructureValidator defines a validator for the commit message structure
type StructureValidator struct{}

// Validate checks the structure of the commit message
func (sv StructureValidator) Validate(c gitlab.Commit) error {
	regex, _ := regexp.Compile("(.*)\n\n([\\S\\s]*)")

	structureMatches := regex.FindStringSubmatch(c.Message)

	if len(structureMatches) == 3 {
		subject := structureMatches[1]
		body := structureMatches[2]

		if subject == "" {
			return errors.New("Subject line cannot be empty")
		}

		if body == "" {
			return errors.New("Body cannot be empty")
		}

		return nil
	}

	return errors.New("Subject-Body structure not present")
}

// SubjectValidator defines a validator for the subject line
type SubjectValidator struct{}

// Validate checks the subject line of the given commit
func (sv SubjectValidator) Validate(c gitlab.Commit) error {
	regex, _ := regexp.Compile("(.*)\n\n")

	subjectMatch := regex.FindStringSubmatch(c.Message)

	if len(subjectMatch) == 2 {
		subject := subjectMatch[1]

		firstChar := subject[:1]

		if strings.ToUpper(firstChar) != firstChar {
			return errors.New("First character of subject must be upper case")
		}

		if len(subject) > 50 {
			return errors.New("Subject line is limited to 50 characters")
		}

		punktuationMarks := []string{"!", "?", "."}

		for _, punktuationMark := range punktuationMarks {
			if strings.HasSuffix(subject, punktuationMark) {
				return errors.New("Subject should not end with a punctuation mark")
			}
		}

		return nil
	}

	return errors.New("Subject-Body structure not present")
}

// BodyValidator defines a validator for the commit body
type BodyValidator struct{}

// Validate checks the body line lengths and structure
func (bv BodyValidator) Validate(c gitlab.Commit) error {
	regex, _ := regexp.Compile("\n\n([\\S\\s]*)")

	bodyMatches := regex.FindStringSubmatch(c.Message)

	if len(bodyMatches) == 2 {
		body := bodyMatches[1]

		bodyLines := strings.Split(body, "\n")

		for _, line := range bodyLines {
			if len(line) > 72 {
				return errors.New("Body lines should be wrapped after 72 characters")
			}
		}

		return nil
	}

	return errors.New("Body cannot be extracted")
}

// TicketNumberValidator defines a validator for ticket numbers in commits
type TicketNumberValidator struct{}

// Validate checks for ticket numbers in the commit message and for their position
func (tnv TicketNumberValidator) Validate(c gitlab.Commit) error {
	regex, _ := regexp.Compile("(.*)\n\n([\\S\\s]*)")

	matches := regex.FindStringSubmatch(c.Message)

	subject := matches[1]
	body := matches[2]

	ticketNumberPattern, _ := regexp.Compile("[\\w]{4,6}-[\\d]{3,6}")

	if ticketNumberPattern.MatchString(subject) {
		return errors.New("Ticket number in subject line is not allowed")
	}

	if ticketNumberPattern.MatchString(body) {
		correctTicketRefs := 0

		prefixedTicketNumberPattern, _ := regexp.Compile("^(Resolves|See):? ([\\w]{4,8}-[\\d]{3,6}(, )*)+$")

		bodyLines := strings.Split(body, "\n")

		for _, line := range bodyLines {
			if prefixedTicketNumberPattern.MatchString(line) {
				correctTicketRefs++
			}
		}

		if correctTicketRefs == 0 {
			return errors.New("Ticket number should be prefixed with either \"Resolves\" or \"See\" in its own line")
		}
	}

	return nil
}
