package model

import (
	"errors"
)

var (
	ErrMissingRepoPathParameter  = errors.New("missing name name path parameter")
	ErrMissingOwnerPathParameter = errors.New("missing owner path parameter")
	ErrAppNotInstalled           = errors.New("GitHub app not installed for given repository")

	ErrIssueMissingAssignee    = errors.New("the issue is missing an assignee")
	ErrIssueMissingClosedAt    = errors.New("the issue is missing the closed at timestamp")
	ErrIssueMissingPullRequest = errors.New("the issue is missing a pull request")

	ErrEventMissingData            = errors.New("the event is missing data promised by the GitHub API")
	ErrEventNotRepoAdded           = errors.New("the event is not a repo added to github event")
	ErrEventNotInstallationCreated = errors.New("the event is not a github created event")

	ErrEventNotHandled = errors.New("the event is not handled")
)
