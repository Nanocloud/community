package apps

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/Nanocloud/community/plaza/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

const domain = "intra.localdomain.com"

type ApplicationParams struct {
	Id             int    `json:"-"`
	CollectionName string `json:"collection-name"`
	Alias          string `json:"alias"`
	DisplayName    string `json:"display-name"`
	FilePath       string `json:"file-path"`
	IconContents   []byte `json:"icon-content"`
}

type ApplicationParamsWin struct {
	Id             int
	CollectionName string
	Alias          string
	DisplayName    string
	FilePath       string
	IconContents   []byte
}

func reterr(e error, resp string, c *echo.Context) error {
	return c.JSON(
		http.StatusInternalServerError,
		hash{
			"error": []hash{
				hash{
					"title":  e.Error(),
					"detail": resp,
				},
			},
		},
	)
}

func retok(c *echo.Context) error {
	return c.JSON(
		http.StatusOK,
		hash{
			"data": hash{
				"success": true,
			},
		},
	)
}

func PublishApp(c *echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err)
		return reterr(err, "", c)
	}
	var p map[string]string
	err = json.Unmarshal(body, &p)
	if err != nil {
		return reterr(err, "", c)
	}
	cmd := exec.Command("powershell.exe", "import-module remotedesktop; New-RDRemoteApp -CollectionName "+p["collection"]+" -DisplayName "+p["displayname"]+" -FilePath "+p["path"])
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return reterr(err, string(resp), c)
	}
	return retok(c)
}

func UnpublishApp(c *echo.Context) error {
	id := c.Param("id")
	username, pwd, _ := c.Request().BasicAuth()
	utils.ExecuteCommandAsAdmin("C:\\Windows\\System32\\WindowsPowershell\\v1.0\\powershell.exe Import-Module RemoteDesktop; Remove-RDRemoteApp -Alias '"+id+"' -CollectionName collection -Force", username, pwd, domain)
	return retok(c)
}

func GetApps(c *echo.Context) error {
	var applications []ApplicationParamsWin
	var winapp ApplicationParamsWin
	var apps []ApplicationParams
	cmd := exec.Command("powershell.exe", "Import-Module RemoteDesktop; Get-RDRemoteApp | ConvertTo-Json -Compress")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return reterr(err, string(resp), c)
	}
	err = json.Unmarshal(resp, &applications)
	if err != nil {
		err = json.Unmarshal(resp, &winapp)
		if err != nil {
			return reterr(err, "", c)
		}
		return c.JSON(
			http.StatusOK,
			ApplicationParams{
				CollectionName: winapp.CollectionName,
				DisplayName:    winapp.DisplayName,
				Alias:          winapp.Alias,
				FilePath:       winapp.FilePath,
				IconContents:   winapp.IconContents,
			},
		)
	}
	for _, app := range applications {
		apps = append(apps, ApplicationParams{
			CollectionName: app.CollectionName,
			DisplayName:    app.DisplayName,
			Alias:          app.Alias,
			FilePath:       app.FilePath,
			IconContents:   app.IconContents,
		})
	}
	return c.JSON(
		http.StatusOK,
		apps,
	)
}
