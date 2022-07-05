package providers

import (
	"context"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// GetReposByOwner lists the repositories that are accessible to an owners authenticated installation.
func (c *githubInstallationClient) GetReposByOwner(ctx context.Context, owner string) ([]string, error) {
	var (
		client, _          = c.clients.get(owner)
		allRepos           []*github.Repository
		allCompressedRepos []string
		listOptions        = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		repoList, resp, err := client.Apps.ListRepos(ctx, listOptions)
		if err != nil {
			return allCompressedRepos, err
		}
		allRepos = append(allRepos, repoList.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		compressedEvent, err := model.NewRepo(repo)
		if err != nil {
			continue
		}
		allCompressedRepos = append(allCompressedRepos, compressedEvent)
	}

	return allCompressedRepos, nil
}
