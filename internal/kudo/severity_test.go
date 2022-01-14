package kudo_test

import (
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"

	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func TestIssueToSeverity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Label    string
		Expected kudo.IssueSeverity
	}{
		{
			Name:     "Issue severity label none",
			Label:    "none",
			Expected: kudo.IssueSeverity("none"),
		},
		{
			Name:     "Issue severity label low",
			Label:    "low",
			Expected: kudo.IssueSeverity("low"),
		},
		{
			Name:     "Issue severity label medium",
			Label:    "medium",
			Expected: kudo.IssueSeverity("medium"),
		},
		{
			Name:     "Issue severity high ",
			Label:    "high",
			Expected: kudo.IssueSeverity("high"),
		},
		{
			Name:     "Issue severity critical ",
			Label:    "critical",
			Expected: kudo.IssueSeverity("critical"),
		},
		{
			Name:     "Issue severity critical ",
			Label:    "",
			Expected: kudo.IssueSeverity("none"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := &github.Issue{Labels: []*github.Label{{Name: &testCase.Label}}}

			// WHEN
			severityResult := kudo.IssueToSeverity(issue)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
		})
	}
}

func TestIssueToSeverityNoIssue(t *testing.T) {
	t.Parallel()
	// GIVEN
	var issue *github.Issue

	// WHEN
	severityResult := kudo.IssueToSeverity(issue)

	// THEN
	assert.Equal(t, kudo.IssueSeverity("none"), severityResult)
}
