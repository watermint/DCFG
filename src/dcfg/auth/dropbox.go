package auth

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"dcfg/config"
	"dcfg/explorer"
	"github.com/cihub/seelog"
	"fmt"
	"os"
	"encoding/json"
)

func verifyDropboxToken(token string) {
	client := dropboxClientFromToken(token)
	team, err := client.GetInfo()
	if err != nil {
		explorer.Fatal("Authentication failed [%s]", err)
	}
	explorer.ReportSuccess("Verified token for Dropbox Team: TeamId[%s] TeamName[%s] Provisioned[%d] Num Licenses[%d]", team.TeamId, team.Name, team.NumProvisionedUsers, team.NumLicensedUsers)
}

func dropboxClientFromToken(token string) dropbox.Api {
	return dropbox.Client(token, dropbox.Options{})
}

func getDropboxTokenFromConsole() string {
	seelog.Flush()

	fmt.Println("Dropbox Business API (permisson type: Team member management)")
	fmt.Println("")
	fmt.Println("------")
	fmt.Println("Paste generated code here:")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		explorer.Fatal("Unable to read authorization code %v", err)
	}

	fmt.Println("")

	return code
}

func updateDropboxToken() {
	token := getDropboxTokenFromConsole()

	verifyDropboxToken(token)

	content := config.DropboxToken{
		TeamManagementToken: token,
	}

	j, err := os.Create(config.Global.DropboxTokenFile())
	if err != nil {
		explorer.Fatal("Unable to open Dropbox token file", config.Global.GoogleTokenFile(), err)
	}
	defer j.Close()

	err = json.NewEncoder(j).Encode(content)
	if err != nil {
		explorer.Fatal("Unable to write Dropbox token file", config.Global.GoogleTokenFile(), err)
	}
	explorer.ReportSuccess("Dropbox Token file updated: [%s]", config.Global.GoogleTokenFile())

}

func DropboxClient() dropbox.Api {
	return dropboxClientFromToken(config.Global.DropboxToken().TeamManagementToken)
}

func AuthDropbox() {
	seelog.Info("Start authentication sequence for Dropbox")
	updateDropboxToken()
}
