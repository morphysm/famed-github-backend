package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

func TestIssue_Severity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Severities  []model.IssueSeverity
		Expected    model.IssueSeverity
		ExpectedErr error
	}{
		{
			Name:       "Severity label none",
			Severities: []model.IssueSeverity{model.Info},
			Expected:   model.Info,
		},
		{
			Name:       "Severity label low",
			Severities: []model.IssueSeverity{model.Low},
			Expected:   model.Low,
		},
		{
			Name:       "Severity label medium",
			Severities: []model.IssueSeverity{model.Medium},
			Expected:   model.Medium,
		},
		{
			Name:       "Severity high ",
			Severities: []model.IssueSeverity{model.High},
			Expected:   model.High,
		},
		{
			Name:       "Severity critical ",
			Severities: []model.IssueSeverity{model.Critical},
			Expected:   model.Critical,
		},
		{
			Name:        "Severity missing",
			Severities:  nil,
			Expected:    "",
			ExpectedErr: model.ErrIssueMissingSeverityLabel,
		},
		{
			Name:        "Multiple severities",
			Severities:  []model.IssueSeverity{model.High, model.Critical},
			Expected:    "",
			ExpectedErr: model.ErrIssueMultipleSeverityLabels,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			issue := model.Issue{Severities: testCase.Severities}

			// WHEN
			severityResult, err := issue.Severity()

			// THEN
			assert.Equal(t, testCase.Expected, severityResult)
			assert.ErrorIs(t, err, testCase.ExpectedErr)
		})
	}
}
