package kudo

import "github.com/google/go-github/v41/github"

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

// IssueToSeverity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns issue severity none
func IssueToSeverity(issue *github.Issue) IssueSeverity {
	if issue == nil || issue.Labels == nil {
		return IssueSeverityNone
	}

	// TODO how do we handle multiple CVSS
	for _, label := range issue.Labels {
		if label.Name == nil {
			continue
		}

		switch *label.Name {
		case string(IssueSeverityLow):
			return IssueSeverityLow
		case string(IssueSeverityMedium):
			return IssueSeverityMedium
		case string(IssueSeverityHigh):
			return IssueSeverityHigh
		case string(IssueSeverityCritical):
			return IssueSeverityCritical
		}
	}

	return IssueSeverityNone
}
