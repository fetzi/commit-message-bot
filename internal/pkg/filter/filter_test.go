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
	FilterTest{"Foo", false},
	FilterTest{"Merge branch foo into bar", true},
	FilterTest{"Revert test", true},
	FilterTest{"Reverts test", false},
	FilterTest{"Automated Jenkins commit", true},
	FilterTest{"This is an automated Jenkins commit", false},
}

var emailFilterTests = []FilterTest{
	FilterTest{"office@karriere.at", false},
	FilterTest{"dev-test@karriere.at", false},
	FilterTest{"lead-dev@karriere.at", false},
	FilterTest{"dev@karriere.at", true},
}

var groupFilterTests = []FilterTest{
	FilterTest{"karriere/site", false},
	FilterTest{"api/api", false},
	FilterTest{"jobs/foo", true},
	FilterTest{"legacy/bar", true},
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
			gitlab.Commit{Message: message},
		},
	}
}

func makeEventWithAuthor(email string) gitlab.PushEventPayload {
	return gitlab.PushEventPayload{
		Commits: []gitlab.Commit{
			gitlab.Commit{
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
