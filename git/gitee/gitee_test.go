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
