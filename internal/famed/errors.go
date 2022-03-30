package famed

import (
	"errors"
)

var (
	ErrIssueMissingAssignee = errors.New("the issue is missing an assignee")
	ErrIssueMissingClosedAt = errors.New("the issue is missing the closed at timestamp")

	ErrEventMissingData            = errors.New("the event is missing data promised by the GitHub API")
	ErrEventNotRepoAdded           = errors.New("the event is not a repo added to github event")
	ErrEventNotInstallationCreated = errors.New("the event is not a github created event")
	ErrEventMissingFamedLabel      = errors.New("the event is missing the famed label")
)
