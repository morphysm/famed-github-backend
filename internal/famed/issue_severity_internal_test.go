package famed

import (
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIssueToSeverity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Label       string
		Expected    config.IssueSeverity
		ExpectedErr error
	}{
		{
			Name:     "Issue severity label none",
			Label:    "none",
			Expected: config.CVSSNone,
		},
		{
			Name:     "Issue severity label low",
			Label:    "low",
			Expected: config.CVSSLow,
		},
		{
			Name:     "Issue severity label medium",
			Label:    "medium",
			Expected: config.CVSSMedium,
		},
		{
			Name:     "Issue severity high ",
			Label:    "high",
			Expected: config.CVSSHigh,
		},
		{
			Name:     "Issue severity critical ",
			Label:    "critical",
			Expected: config.CVSSCritical,
		},
		{
			Name:        "Issue severity critical ",
			Label:       "",
			Expected:    "",
			ExpectedErr: ErrIssueMissingSeverityLabel,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := Issue{Issue: &github.Issue{Labels: []*github.Label{{Name: &testCase.Label}}}}

			// WHEN
			severityResult, err := issue.severity()

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
	labelNone := string(config.CVSSNone)
	labelLow := string(config.CVSSCritical)
	issue := Issue{Issue: &github.Issue{Labels: []*github.Label{{Name: &labelNone}, {Name: &labelLow}}}}

	// WHEN
	severityResult, err := issue.severity()

	// THEN
	assert.Equal(t, config.IssueSeverity(""), severityResult)
	assert.ErrorIs(t, ErrIssueMultipleSeverityLabels, err)
}
