package sessions

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func formatResponse(tab []string, id string) [][]string {
	var format [][]string
	for _, val := range tab {
		newtab := strings.Fields(val)
		if len(newtab) == 4 {
			if id == "Administrator" || newtab[1] == id {
				format = append(format, newtab)
			}
		}
	}
	return format
}

func Get(c *echo.Context) error {
	cmd := exec.Command("powershell.exe", "query session | ConvertTo-Json -Compress")
	resp, err := cmd.CombinedOutput()
	var tab []string
	err = json.Unmarshal(resp, &tab)
	if err != nil {
		log.Error("Error while unmarshaling query response: ", err)
	}
	response := formatResponse(tab, c.Param("id"))
	if len(response) == 0 {
		response = make([][]string, 0)
	}
	return c.JSON(
		http.StatusOK,
		hash{
			"data": response,
		},
	)
}

func Logoff(c *echo.Context) error {
	cmd := exec.Command("powershell.exe", "query session | ConvertTo-Json -Compress")
	resp, err := cmd.CombinedOutput()
	var tab []string
	err = json.Unmarshal(resp, &tab)
	if err != nil {
		log.Error("Error while unmarshaling query response: ", err)
	}
	response := formatResponse(tab, c.Param("id"))
	if len(response) == 1 {
		cmd := exec.Command("powershell.exe", "logoff "+response[0][2])
		resp, err = cmd.CombinedOutput()
		if err != nil {
			log.Error("Error while loging off user: ", err, string(resp))
		}
	}
	return c.JSON(
		http.StatusOK,
		hash{
			"data": response,
		},
	)
}
