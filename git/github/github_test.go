package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestRepositories(t *testing.T) {
	accessToken := os.Getenv("GITHUB_PNCX_TOKEN")
	user := "pncx"
	repos, err := Repositories(user, accessToken)
	fmt.Println("repos", repos, "err", err)
	for _, repo := range repos {
		fmt.Println("index:", repo.GetID(), "org:", repo.GetOrganization(), "name:", repo.GetName(),
			"FullName:", repo.GetFullName(), "CloneURL:", repo.GetCloneURL(), "SSHURL:", repo.GetSSHURL())
	}
}

func TestRepositoriesByOrg(t *testing.T) {
	accessToken := os.Getenv("GITHUB_PNCX_TOKEN")
	org := "x-cx"
	repos, err := RepositoriesByOrg(org, accessToken, true)
	fmt.Println("repos", repos, "err", err)
	for _, repo := range repos {
		fmt.Println("id:", repo.GetID(), "org:", repo.GetOrganization(), "name:", repo.GetName(),
			"FullName:", repo.GetFullName(), "CloneURL:", repo.GetCloneURL(), "SSHURL:", repo.GetSSHURL())
	}
}

func TestOrganizations(t *testing.T) {
	accessToken := os.Getenv("GITHUB_PNCX_TOKEN")
	user := ""
	orgs, err := Organizations(user, accessToken)
	fmt.Println(orgs, err)
	for _, org := range orgs {
		fmt.Println("id:", org.GetID(), "org:", "Login:", org.GetLogin())
		repos, err := RepositoriesByOrg(org.GetLogin(), accessToken, true)
		fmt.Println("repos", repos, "err", err)
		reposJosn, err := json.Marshal(repos)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(reposJosn))

		for _, repo := range repos {
			if repo.GetPermissions()["admin"] {
				fmt.Println("  id:", repo.GetID(), "org:", repo.GetOrganization(), "name:", repo.GetName(),
					"FullName:", repo.GetFullName(), "CloneURL:", repo.GetCloneURL(), "SSHURL:", repo.GetSSHURL())
			} else {
				fmt.Println("  not admin", org.GetLogin())
			}
		}
	}
}
