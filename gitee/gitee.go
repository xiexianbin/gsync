package gitee

import (
	"context"
	"github.com/antihax/optional"
	"golang.org/x/oauth2"

	"gitee.com/openeuler/go-gitee/gitee"
)

func giteeClient(accessToken string) (*gitee.APIClient, context.Context) {
	// oauth
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	// configuration
	giteeConf := gitee.NewConfiguration()
	giteeConf.HTTPClient = oauth2.NewClient(ctx, ts)
	// git client
	giteeClient := gitee.NewAPIClient(giteeConf)

	return giteeClient, ctx
}

// list Organizations
func Organizations(user, accessToken string) ([]gitee.Group, error){
	// gitee client
	giteeClient, ctx := giteeClient(accessToken)
	opt := &gitee.GetV5UsersUsernameOrgsOpts{
		AccessToken: optional.NewString(accessToken),
		PerPage: optional.NewInt32(1000),
	}
	groups, _, err := giteeClient.OrganizationsApi.GetV5UsersUsernameOrgs(ctx, user, opt)

	return groups, err
}

// return user Repositories
func Repositories(accessToken string, page, perPage int) (gitee.Project, error){
	// gitee client
	giteeClient, ctx := giteeClient(accessToken)

	// list all repositories for the authenticated user
	opt := &gitee.GetV5UserReposOpts{
		AccessToken: optional.NewString(accessToken),
		Page: optional.NewInt32(int32(page)),
		PerPage: optional.NewInt32(int32(perPage)),
	}
	project, _, err := giteeClient.RepositoriesApi.GetV5UserRepos(ctx, opt)

	return project, err
}

// return user RepositoriesByOrg
func RepositoriesByOrg(org, accessToken string, page, perPage int) ([]gitee.Project, error){
	// gitee client
	giteeClient, ctx := giteeClient(accessToken)

	// list repositories for org "gitee"
	opt := &gitee.GetV5OrgsOrgReposOpts{
		AccessToken: optional.NewString(accessToken),
		Page: optional.NewInt32(int32(page)),
		PerPage: optional.NewInt32(int32(perPage)),
	}
	projects, _, err := giteeClient.RepositoriesApi.GetV5OrgsOrgRepos(ctx, org, opt)

	return projects, err
}

func CreateUserRepos(accessToken, name, description, homepage string, private, hasIssues, hasWiki, canComment, autoInit bool) (gitee.Project, error){
	// gitee client
	giteeClient, ctx := giteeClient(accessToken)

	// list repositories for org "gitee"
	opt := &gitee.PostV5UserReposOpts{
		AccessToken: optional.NewString(accessToken),
	}

	if description != "" {
		opt.Description = optional.NewString(description)
	}
	if homepage != "" {
		opt.Homepage = optional.NewString(homepage)
	}
	if private {
		opt.Private = optional.NewBool(private)
	}
	if hasIssues {
		opt.HasIssues = optional.NewBool(hasIssues)
	}
	if hasWiki {
		opt.HasWiki = optional.NewBool(hasWiki)
	}
	if canComment {
		opt.CanComment = optional.NewBool(canComment)
	}
	if autoInit {
		opt.AutoInit = optional.NewBool(autoInit)
	}

	project, _, err := giteeClient.RepositoriesApi.PostV5UserRepos(ctx, name, opt)

	return project, err
}

