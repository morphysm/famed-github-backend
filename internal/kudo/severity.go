package kudo

import (
	"errors"

	"github.com/google/go-github/v41/github"
)

type IssueSeverity string

const (
	// IssueSeverityNone represents a CVSS of 0
	IssueSeverityNone IssueSeverity = "none"
	// IssueSeverityLow represents a CVSS of 0.1-3.9
	IssueSeverityLow IssueSeverity = "low"
	// IssueSeverityMedium represents a CVSS of 4.0-6.9
	IssueSeverityMedium IssueSeverity = "medium"
	// IssueSeverityHigh represents a CVSS of 7.0-8.9
	IssueSeverityHigh IssueSeverity = "high"
	// IssueSeverityCritical represents a CVSS of 9.0-10.0
	IssueSeverityCritical IssueSeverity = "critical"
)

var (
	ErrIssueMissingSeverityLabel   = errors.New("the issue is missing it's severity label")
	ErrIssueMultipleSeverityLabels = errors.New("the issue has multiple severity labels")
)

// IssueToSeverity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns the IssueMissingLabelErr
// if multiple matching issue severity labels can be found it returns the IssueMultipleSeverityLabelsErr.
func IssueToSeverity(issue *github.Issue) (IssueSeverity, error) {
	var severity IssueSeverity
	for _, label := range issue.Labels {
		if !isLabelValid(label) {
			continue
		}

		// Check if label is equal to one of the predefined severity labels.
		if *label.Name == string(IssueSeverityNone) ||
			*label.Name == string(IssueSeverityLow) ||
			*label.Name == string(IssueSeverityMedium) ||
			*label.Name == string(IssueSeverityHigh) ||
			*label.Name == string(IssueSeverityCritical) {
			// If
			if severity != "" {
				return "", ErrIssueMultipleSeverityLabels
			}
			severity = IssueSeverity(*label.Name)
		}
	}

	if severity == "" {
		return "", ErrIssueMissingSeverityLabel
	}

	return severity, nil
}
