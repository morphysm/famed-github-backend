package github

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type Repo struct {
	Name string
}

func (c *githubInstallationClient) GetRepos(ctx context.Context, owner string) ([]Repo, error) {
	var (
		client, _          = c.clients.get(owner)
		allRepos           []*github.Repository
		allCompressedRepos []Repo
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
		compressedEvent, err := validateRepo(repo)
		if err != nil {
			continue
		}
		allCompressedRepos = append(allCompressedRepos, compressedEvent)
	}

	return allCompressedRepos, nil
}

func validateRepo(repo *github.Repository) (Repo, error) {
	if repo == nil ||
		repo.Name == nil {
		return Repo{}, ErrRepoMissingData
	}

	return Repo{
		Name: *repo.Name,
	}, nil
}
