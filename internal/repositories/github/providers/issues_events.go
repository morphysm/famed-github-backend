package providers

import (
	"github.com/phuslu/log"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/pkg/parse"
)

func (c *githubInstallationClient) ValidateWebHookEvent(request *http.Request) (interface{}, error) {
	var event interface{}
	webhookSecret := []byte(c.webhookSecret)

	payload, err := github.ValidatePayload(request, webhookSecret)
	if err != nil {
		return nil, err
	}

	event, err = github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return nil, err
	}

	switch event := event.(type) {
	case *github.IssuesEvent:
		issuesEvent, err := model.NewIssuesEvent(event, c.famedLabel)
		if err != nil {
			return nil, err
		}

		// TODO see if this code duplication can be cleaned up
		if issuesEvent.Issue.Migrated {
			// Parse red team from issue body
			redTeam, err := parse.FindRightOfKey(*event.Issue.Body, "Bounty Hunter:")
			if err != nil {
				return nil, err
			}

			// Split bounty hunters if two are present separated by ", "
			splitTeam := strings.Split(redTeam, ", ")

			for _, pseudonym := range splitTeam {
				redTeamer, err := c.getRedTeamer(request.Context(), *event.Repo.Owner.Login, pseudonym)
				if err != nil {
					return nil, err
				}
				issuesEvent.Issue.RedTeam = append(issuesEvent.Issue.RedTeam, redTeamer)
			}
		}

		return issuesEvent, err
	case *github.InstallationRepositoriesEvent:
		installationRepositoriesEvent, err := model.NewInstallationRepositoriesEvent(event)
		if err != nil {
			return nil, err
		}

		return installationRepositoriesEvent, err
	case *github.InstallationEvent:
		installationEvent, err := model.NewInstallationEvent(event)
		if err != nil {
			return nil, err
		}

		return installationEvent, err
	default:
		log.Error().Msg("[ValidateWebHookEvent] unhandled event")
		return event, model.ErrUnhandledEventType
	}
}
