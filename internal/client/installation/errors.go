package installation

import "errors"

var (
	ErrUnhandledEventType      = errors.New("the event is not handled")
	ErrIssueMissingData        = errors.New("the issue is missing data promised by the GitHub API")
	ErrIssueCommentMissingData = errors.New("the issue comment is missing data promised by the GitHub API")
	ErrEventMissingData        = errors.New("the event is missing data promised by the GitHub API")
)
