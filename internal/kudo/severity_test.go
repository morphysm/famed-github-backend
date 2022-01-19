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
		Name        string
		Label       string
		Expected    kudo.IssueSeverity
		ExpectedErr error
	}{
		{
			Name:     "Issue severity label none",
			Label:    "none",
			Expected: kudo.IssueSeverityNone,
		},
		{
			Name:     "Issue severity label low",
			Label:    "low",
			Expected: kudo.IssueSeverityLow,
		},
		{
			Name:     "Issue severity label medium",
			Label:    "medium",
			Expected: kudo.IssueSeverityMedium,
		},
		{
			Name:     "Issue severity high ",
			Label:    "high",
			Expected: kudo.IssueSeverityHigh,
		},
		{
			Name:     "Issue severity critical ",
			Label:    "critical",
			Expected: kudo.IssueSeverityCritical,
		},
		{
			Name:        "Issue severity critical ",
			Label:       "",
			Expected:    "",
			ExpectedErr: kudo.ErrIssueMissingLabel,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := &github.Issue{Labels: []*github.Label{{Name: &testCase.Label}}}

			// WHEN
			severityResult, err := kudo.IssueToSeverity(issue)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
			if testCase.ExpectedErr != nil {
				assert.ErrorIs(t, testCase.ExpectedErr, err)
			}
		})
	}
}

func TestIssueToSeverityMultipleSeverityLabels(t *testing.T) {
	t.Parallel()
	// GIVEN
	labelNone := string(kudo.IssueSeverityNone)
	labelLow := string(kudo.IssueSeverityCritical)
	issue := &github.Issue{Labels: []*github.Label{{Name: &labelNone}, {Name: &labelLow}}}

	// WHEN
	severityResult, err := kudo.IssueToSeverity(issue)

	// THEN
	assert.Equal(t, kudo.IssueSeverity(""), severityResult)
	assert.ErrorIs(t, kudo.ErrIssueMultipleSeverityLabels, err)
}
