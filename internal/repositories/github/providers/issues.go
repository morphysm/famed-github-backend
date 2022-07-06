package providers

import (
	"context"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/phuslu/log"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/pkg/parse"
)

// IssueListByRepoOptions specifies filtering options to be passed
// to the GetIssuesByRepo method.
type IssueListByRepoOptions struct {
	Labels   []string
	State    *model.IssueState
	Assignee *string
}

// GetIssuesByRepo returns all issues from a given repository.
func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, owner string, repoName string, options IssueListByRepoOptions) ([]model.Issue, error) {
	var (
		client, _           = c.clients.get(owner)
		allIssues           []*github.Issue
		allCompressedIssues []model.Issue
		listOptions         = &github.IssueListByRepoOptions{
			Labels: options.Labels,
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 30,
			},
		}
	)

	if options.Assignee != nil {
		listOptions.Assignee = *options.Assignee
	}
	if options.State != nil {
		listOptions.State = string(*options.State)
	} else {
		listOptions.State = string(model.All)
	}

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, listOptions)
		if err != nil {
			return allCompressedIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, issue := range allIssues {
		compressedIssue, err := model.NewIssue(issue, owner, repoName)
		if err != nil {
			log.Error().Err(err).Msgf("[GetIssuesByRepo] validation error for issue with number %d", issue.Number)
		}

		if compressedIssue.Migrated {
			// Parse red team from issue body
			redTeam, err := parse.FindRightOfKey(*issue.Body, "Bounty Hunter:")
			if err != nil {
				return nil, err
			}

			// Split bounty hunters if two are present separated by ", "
			splitTeam := strings.Split(redTeam, ", ")

			for _, pseudonym := range splitTeam {
				redTeamer, err := c.getRedTeamer(ctx, owner, pseudonym)
				if err != nil {
					return nil, err
				}
				compressedIssue.RedTeam = append(compressedIssue.RedTeam, redTeamer)
			}
		}

		allCompressedIssues = append(allCompressedIssues, compressedIssue)
	}

	return allCompressedIssues, nil
}
