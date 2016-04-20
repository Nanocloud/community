package exec

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Nanocloud/community/plaza/windows"
	"github.com/labstack/echo"
)

type bodyRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Domain   string   `json:"domain"`
	Command  []string `json:"command"`
	Stdin    string   `json:"stdin"`
}

func Route(c *echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	body := bodyRequest{}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}

	cmd := windows.Command(
		body.Username, body.Domain, body.Password,
		body.Command[0], body.Command[1:]...,
	)

	if body.Stdin != "" {
		cmd.Stdin = strings.NewReader(body.Stdin)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, err.Error())
	// }

	res := make(map[string]interface{})
	res["stdout"] = stdout.String()
	res["stderr"] = stderr.String()
	res["time"] = cmd.ProcessState.SysUsage()
	res["code"] = cmd.ProcessState.Status.ExitCode

	return c.JSON(http.StatusOK, res)
}
