package filter

import (
	"testing"

	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

type FilterTest struct {
	value  string
	result bool
}

var bodyStartsWithTests = []FilterTest{
	{"Foo", false},
	{"Merge branch foo into bar", true},
	{"Revert test", true},
	{"Reverts test", false},
	{"Automated Jenkins commit", true},
	{"This is an automated Jenkins commit", false},
}

var emailFilterTests = []FilterTest{
	{"office@karriere.at", false},
	{"dev-test@karriere.at", false},
	{"lead-dev@karriere.at", false},
	{"dev@karriere.at", true},
}

var groupFilterTests = []FilterTest{
	{"karriere/site", false},
	{"api/api", false},
	{"jobs/foo", true},
	{"legacy/bar", true},
}

func TestBodyStartsWithFilter(t *testing.T) {
	filter := BodyStartsWithFilter{
		[]string{"Merge branch", "Merge remote-tracking", "Revert ", "Automated Jenkins commit"},
	}

	for _, test := range bodyStartsWithTests {
		event := makeEventWithCommit(test.value)
		result := filter.Filter(event)

		if result != test.result {
			t.Errorf(
				"BodyStartsWithFilter for \"%s\" failed. Expected %t but got %t",
				test.value,
				test.result,
				result,
			)
		}
	}
}

func TestEmailFilter(t *testing.T) {
	filter := EmailFilter{
		[]string{"dev@karriere.at"},
	}

	for _, test := range emailFilterTests {
		event := makeEventWithAuthor(test.value)
		result := filter.Filter(event)

		if result != test.result {
			t.Errorf(
				"EmailFilter for \"%s\" failed. Expected %t but got %t",
				test.value,
				test.result,
				result,
			)
		}
	}
}

func TestGroupFilter(t *testing.T) {
	filter := GroupFilter{
		[]string{"jobs", "internal", "legacy"},
	}

	for _, test := range groupFilterTests {
		event := makeEventWithProject(test.value)
		result := filter.Filter(event)

		if result != test.result {
			t.Errorf(
				"GroupFilter for \"%s\" failed. Expected %t but got %t",
				test.value,
				test.result,
				result,
			)
		}
	}
}

func makeEventWithCommit(message string) gitlab.PushEventPayload {
	return gitlab.PushEventPayload{
		Commits: []gitlab.Commit{
			{Message: message},
		},
	}
}

func makeEventWithAuthor(email string) gitlab.PushEventPayload {
	return gitlab.PushEventPayload{
		Commits: []gitlab.Commit{
			{
				Author: gitlab.Author{Email: email},
			},
		},
	}
}

func makeEventWithProject(namespace string) gitlab.PushEventPayload {
	return gitlab.PushEventPayload{
		Project: gitlab.Project{PathWithNamespace: namespace},
	}
}
