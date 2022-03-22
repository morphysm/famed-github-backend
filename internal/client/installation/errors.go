package installation

import "errors"

var (
	//ErrIssueMissingAssignee   = errors.New("the issue is missing an assignee")
	ErrIssueMissingData = errors.New("the issue is missing data promised by the GitHub API")
	//ErrIssueMissingFamedLabel = errors.New("the issue is missing the famed label")
	//
	ErrUnhandledEventType = errors.New("the event is not handled")
	ErrEventMissingData   = errors.New("the event is missing data promised by the GitHub API")
	//ErrEventAssigneeMissingData    = errors.New("the event assignee is missing data promised by the GitHub API")
	//ErrEventNotClose               = errors.New("the event is not a close event")
	//ErrEventNotRepoAdded           = errors.New("the event is not a repo added to installation event")
	//ErrEventNotInstallationCreated = errors.New("the event is not a installation created event")
	//ErrEventMissingFamedLabel      = errors.New("the event is missing the famed label")
)
