package gitee

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestOrganizations(t *testing.T) {
	accessToken := os.Getenv("GITEE_PNCX_TOKEN")
	user := "pncx"
	groups, err := Organizations(user, accessToken)
	fmt.Println("groups", groups, "err", err)
	reposJosn, err := json.Marshal(groups)
	if err != nil{
		log.Fatalln(err)
	}

	fmt.Println(string(reposJosn))
	for _, group := range groups {
		fmt.Println("id:", group.Id, "org:", group.Login, "url:", group.Url)
	}
}

func TestRepositories(t *testing.T) {
	accessToken := os.Getenv("GITEE_PNCX_TOKEN")
	repos, err := Repositories(accessToken, 1, 1000)
	fmt.Println("repos", repos, "err", err)
	reposJosn, err := json.Marshal(repos)
	if err != nil{
		log.Fatalln(err)
	}

	fmt.Println(string(reposJosn))
	//for _, repo := range repos {
	//	fmt.Println("id:", group.Id, "org:", group.Login, "url:", group.Url)
	//}
}

func TestRepositoriesByOrg(t *testing.T) {
	accessToken := os.Getenv("GITEE_PNCX_TOKEN")
	org := "openeuler"
	repos, err := RepositoriesByOrg(org, accessToken, 1, 1000)
	fmt.Println("repos", repos, "err", err)
	reposJosn, err := json.Marshal(repos)
	if err != nil{
		log.Fatalln(err)
	}

	fmt.Println(string(reposJosn))
	for _, repo := range repos {
		fmt.Println("id:", repo.Id, "name:", repo.Name, "url:", repo.Url)
	}
}

func TestCreateUserRepos(t *testing.T) {
	accessToken := os.Getenv("GITEE_PNCX_TOKEN")
	project, err := CreateUserRepos(accessToken, "abc", "", "", false, false, false, false, false)
	fmt.Println("project", project, "err", err)
	projectJosn, err := json.Marshal(project)
	if err != nil{
		log.Fatalln(err)
	}

	fmt.Println(string(projectJosn))
}
