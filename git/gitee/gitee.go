package gitee

import (
	"context"
	"github.com/antihax/optional"
	"golang.org/x/oauth2"

	"gitee.com/openeuler/go-gitee/gitee"
)

// list Organizations
func Organizations(user, accessToken string) ([]gitee.Group, error){
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
	opt := &gitee.GetV5UsersUsernameOrgsOpts{
		AccessToken: optional.NewString(accessToken),
		PerPage: optional.NewInt32(1000),
	}
	groups, _, err := giteeClient.OrganizationsApi.GetV5UsersUsernameOrgs(ctx, user, opt)

	return groups, err
}
