package model

import "github.com/google/go-github/v41/github"

type IssueSeverity string

const (
	// Info represents a CVSS of 0
	Info IssueSeverity = "info"
	// Low represents a CVSS of 0.1-3.9
	Low IssueSeverity = "low"
	// Medium represents a CVSS of 4.0-6.9
	Medium IssueSeverity = "medium"
	// High represents a CVSS of 7.0-8.9
	High IssueSeverity = "high"
	// Critical represents a CVSS of 9.0-10.0
	Critical IssueSeverity = "critical"
)

// newSeverity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns the IssueMissingLabelErr
// if multiple matching issue severity labels can be found it returns the IssueMultipleSeverityLabelsErr.
func newSeverity(labels []*github.Label) []IssueSeverity {
	var severities []IssueSeverity
	for _, label := range labels {
		// Check if label is equal to one of the predefined severity labels.
		if (label != nil &&
			label.Name != nil) &&
			(*label.Name == string(Info) ||
				*label.Name == string(Low) ||
				*label.Name == string(Medium) ||
				*label.Name == string(High) ||
				*label.Name == string(Critical)) {
			severities = append(severities, IssueSeverity(*label.Name))
		}
	}

	return severities
}
