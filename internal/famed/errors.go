package famed

import (
	"errors"
)

var (
	ErrIssueMissingAssignee = errors.New("the issue is missing an assignee")
	ErrIssueClosedAt        = errors.New("the issue is missing the closed at timestamp")

	ErrEventMissingData            = errors.New("the event is missing data promised by the GitHub API")
	ErrEventNotRepoAdded           = errors.New("the event is not a repo added to installation event")
	ErrEventNotInstallationCreated = errors.New("the event is not a installation created event")
	ErrEventMissingFamedLabel      = errors.New("the event is missing the famed label")
)
