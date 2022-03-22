package famed

import (
	"testing"

	"github.com/morphysm/famed-github-backend/internal/client/installation"
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
			Name:     "WrappedIssue severity label none",
			Label:    "none",
			Expected: config.CVSSNone,
		},
		{
			Name:     "WrappedIssue severity label low",
			Label:    "low",
			Expected: config.CVSSLow,
		},
		{
			Name:     "WrappedIssue severity label medium",
			Label:    "medium",
			Expected: config.CVSSMedium,
		},
		{
			Name:     "WrappedIssue severity high ",
			Label:    "high",
			Expected: config.CVSSHigh,
		},
		{
			Name:     "WrappedIssue severity critical ",
			Label:    "critical",
			Expected: config.CVSSCritical,
		},
		{
			Name:        "WrappedIssue severity critical ",
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
			issue := WrappedIssue{Issue: installation.Issue{Labels: []installation.Label{{Name: testCase.Label}}}}

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
	issue := WrappedIssue{Issue: installation.Issue{Labels: []installation.Label{{Name: labelNone}, {Name: labelLow}}}}

	// WHEN
	severityResult, err := issue.severity()

	// THEN
	assert.Equal(t, config.IssueSeverity(""), severityResult)
	assert.ErrorIs(t, ErrIssueMultipleSeverityLabels, err)
}
