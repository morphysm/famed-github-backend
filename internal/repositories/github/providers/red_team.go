package providers

import (
	"context"
	"log"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// parseRedTeam returns a red team parsed from a GitHub issue body.
func (c *githubInstallationClient) getRedTeamer(ctx context.Context, owner string, pseudonym string) (model.User, error) {
	// Read from known GitHub logins
	login := c.redTeamLogins[pseudonym]
	if login == "" {
		log.Printf("[parseRedTeam] no GitHub login found for red teamer %s", pseudonym)
		return model.User{Login: pseudonym}, nil
	}

	// Check if red teamer is in cache
	cachedTeamer, ok := c.cachedRedTeam.Get(login)
	if ok {
		return cachedTeamer, nil
	}

	// Fetch user info
	redTeamer, err := c.GetUser(ctx, owner, login)
	if err != nil {
		return model.User{}, err
	}

	// Add user info to cache
	c.cachedRedTeam.Add(redTeamer)
	return redTeamer, nil
}
