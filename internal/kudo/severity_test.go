package kudo

import (
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestIssueToSeverity(t *testing.T) {
	testCases := []struct {
		Name     string
		Label    string
		Expected IssueSeverity
	}{
		{
			Name:     "Issue severity label none",
			Label:    "none",
			Expected: IssueSeverity("none"),
		},
		{
			Name:     "Issue severity label low",
			Label:    "low",
			Expected: IssueSeverity("low"),
		},
		{
			Name:     "Issue severity label medium",
			Label:    "medium",
			Expected: IssueSeverity("medium"),
		},
		{
			Name:     "Issue severity high ",
			Label:    "high",
			Expected: IssueSeverity("high"),
		},
		{
			Name:     "Issue severity critical ",
			Label:    "critical",
			Expected: IssueSeverity("critical"),
		},
		{
			Name:     "Issue severity critical ",
			Label:    "",
			Expected: IssueSeverity("none"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := &github.Issue{Labels: []*github.Label{{Name: &testCase.Label}}}

			// WHEN
			severityResult := IssueToSeverity(issue)

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
	severityResult := IssueToSeverity(issue)

	// THEN
	assert.Equal(t, IssueSeverity("none"), severityResult)
}
