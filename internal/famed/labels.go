package famed

import (
	"errors"
	"log"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
)

var (
	ErrIssueMissingSeverityLabel   = errors.New("the issue is missing it's severity label")
	ErrIssueMultipleSeverityLabels = errors.New("the issue has multiple severity labels")
)

// severity returns the issue severity by matching labels against CVSS
// if no matching issue severity label can be found it returns the IssueMissingLabelErr
// if multiple matching issue severity labels can be found it returns the IssueMultipleSeverityLabelsErr.
func severity(issue github.Issue) (config.IssueSeverity, error) {
	var severity config.IssueSeverity
	for _, label := range issue.Labels {
		// Check if label is equal to one of the predefined severity labels.
		if label.Name == string(config.CVSSNone) ||
			label.Name == string(config.CVSSLow) ||
			label.Name == string(config.CVSSMedium) ||
			label.Name == string(config.CVSSHigh) ||
			label.Name == string(config.CVSSCritical) {
			// If
			if severity != "" {
				return "", ErrIssueMultipleSeverityLabels
			}
			severity = config.IssueSeverity(label.Name)
		}
	}

	if severity == "" {
		return "", ErrIssueMissingSeverityLabel
	}

	return severity, nil
}

// isIssueFamedLabeled checks weather the issue labels contain expected famed label.
func isIssueFamedLabeled(issue github.Issue, famedLabel string) bool {
	for _, label := range issue.Labels {
		if label.Name == famedLabel {
			return true
		}
	}

	log.Printf("[IsIssueFamedLabeled] missing famed label: %s in issue with ID: %d", famedLabel, issue.ID)
	return false
}
