package github

import "errors"

var (
	ErrUnhandledEventType      = errors.New("the event is not handled")
	ErrInstallationMissingData = errors.New("the installation is missing data promised by the GitHub API")
	ErrRepoMissingData         = errors.New("the repo is missing data promised by the GitHub API")
	ErrIssueMissingData        = errors.New("the issue is missing data promised by the GitHub API")
	ErrUserMissingData         = errors.New("the user is missing data promised by the GitHub API")
	ErrIssueCommentMissingData = errors.New("the issue comment is missing data promised by the GitHub API")
	ErrEventMissingData        = errors.New("the event is missing data promised by the GitHub API")
)
