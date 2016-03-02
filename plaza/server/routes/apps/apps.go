package apps

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type hash map[string]interface{}

var nanocloud = ""

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

func Nano(c *echo.Context) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		reterr(err, "", c)
	}
	var p map[string]string
	err = json.Unmarshal(body, &p)
	if err != nil {
		reterr(err, "", c)
	}
	nanocloud = p["address"]
}

func PublishApp(c *echo.Context) {
	if nanocloud == "" {
		reterr(errors.New("Nanocloud address not provided"), "", c)
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		reterr(err, "", c)
	}
	var p map[string]string
	err = json.Unmarshal(body, &p)
	if err != nil {
		reterr(err, "", c)
	}
	cmd := exec.Command("powershell.exe", "import-module remotedesktop; New-RDRemoteApp -CollectionName "+p["collection"]+" -DisplayName "+p["displayname"]+" -FilePath "+p["path"])
	resp, err := cmd.CombinedOutput()
	if err != nil {
		reterr(err, string(resp), c)
	}
}
