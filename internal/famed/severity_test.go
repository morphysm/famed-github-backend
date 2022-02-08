package famed_test

import (
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/stretchr/testify/assert"
)

func TestIssueToSeverity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Label       string
		Expected    famed.IssueSeverity
		ExpectedErr error
	}{
		{
			Name:     "Issue severity label none",
			Label:    "none",
			Expected: famed.IssueSeverityNone,
		},
		{
			Name:     "Issue severity label low",
			Label:    "low",
			Expected: famed.IssueSeverityLow,
		},
		{
			Name:     "Issue severity label medium",
			Label:    "medium",
			Expected: famed.IssueSeverityMedium,
		},
		{
			Name:     "Issue severity high ",
			Label:    "high",
			Expected: famed.IssueSeverityHigh,
		},
		{
			Name:     "Issue severity critical ",
			Label:    "critical",
			Expected: famed.IssueSeverityCritical,
		},
		{
			Name:        "Issue severity critical ",
			Label:       "",
			Expected:    "",
			ExpectedErr: famed.ErrIssueMissingSeverityLabel,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := &github.Issue{Labels: []*github.Label{{Name: &testCase.Label}}}

			// WHEN
			severityResult, err := famed.IssueToSeverity(issue)

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
	labelNone := string(famed.IssueSeverityNone)
	labelLow := string(famed.IssueSeverityCritical)
	issue := &github.Issue{Labels: []*github.Label{{Name: &labelNone}, {Name: &labelLow}}}

	// WHEN
	severityResult, err := famed.IssueToSeverity(issue)

	// THEN
	assert.Equal(t, famed.IssueSeverity(""), severityResult)
	assert.ErrorIs(t, famed.ErrIssueMultipleSeverityLabels, err)
}
