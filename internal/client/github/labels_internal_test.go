package github

import (
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

func TestIssueToSeverity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Label    string
		Expected []IssueSeverity
	}{
		{
			Name:     "Severity label none",
			Label:    "info",
			Expected: []IssueSeverity{Info},
		},
		{
			Name:     "Severity label low",
			Label:    "low",
			Expected: []IssueSeverity{Low},
		},
		{
			Name:     "Severity label medium",
			Label:    "medium",
			Expected: []IssueSeverity{Medium},
		},
		{
			Name:     "Severity high ",
			Label:    "high",
			Expected: []IssueSeverity{High},
		},
		{
			Name:     "Severity critical ",
			Label:    "critical",
			Expected: []IssueSeverity{Critical},
		},
		{
			Name:     "Severity missing",
			Label:    "",
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			labels := []*github.Label{{Name: pointer.String(testCase.Label)}}

			// WHEN
			severityResult := parseSeverities(labels)

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
		})
	}
}

func TestIssueToSeverityMultipleSeverityLabels(t *testing.T) {
	t.Parallel()
	// GIVEN
	labelNone := string(Info)
	labelLow := string(Critical)
	labels := []*github.Label{{Name: pointer.String(labelNone)}, {Name: pointer.String(labelLow)}}

	// WHEN
	severityResult := parseSeverities(labels)

	// THEN
	assert.Equal(t, []IssueSeverity{Info, Critical}, severityResult)
}
