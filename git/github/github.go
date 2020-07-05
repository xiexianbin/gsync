package github

import (
	"context"
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

// list Organizations
func Organizations(user, accessToken string) ([]*github.Organization, error){
	ctx := context.Background()
	client := github.NewClient(nil)
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	opt := &github.ListOptions{
		PerPage: 1000,
	}
	orgs, _, err := client.Organizations.List(ctx, user, opt)

	return orgs, err
}

// return user Repositories
func Repositories(user, accessToken string) ([]*github.Repository, error){
	ctx := context.Background()
	client := github.NewClient(nil)
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, user, nil)

	return repos, err
}

// return user RepositoriesByOrg
func RepositoriesByOrg(org, accessToken string, private bool) ([]*github.Repository, error){
	ctx := context.Background()
	client := github.NewClient(nil)
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	// list repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{}
	if private == false {
		opt.Type = "public" // Can be one of public, private or internal
	}
	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

	return repos, err
}
