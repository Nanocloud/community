package sessions

import (
	"io/ioutil"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/models/sessions"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

var kServer string
var kPort string

type hash map[string]interface{}

func List(c *echo.Context) error {

	user := c.Get("user").(*users.User)

	winUser, err := user.WindowsCredentials()
	if err != nil {
		return err
	}

	sessionList, err := sessions.GetAll(winUser.Sam)

	if err != nil {
		log.Error(err)
		return utils.JSON(c, http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}

	var response = make([]hash, len(sessionList))
	for i, val := range sessionList {
		res := hash{
			"id":         val.Id,
			"type":       "session",
			"attributes": val,
		}
		response[i] = res
	}

	return c.JSON(http.StatusOK, hash{"data": response})
}

func Logoff(c *echo.Context) error {
	user := c.Get("user").(*users.User)

	winUser, err := user.WindowsCredentials()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", "http://"+kServer+":"+kPort+"/sessions/"+winUser.Sam, nil)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.Status != "200 OK" {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": resp.Status,
				},
			},
		})
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}
	return c.JSON(http.StatusOK, string(b))
}

func init() {
	kServer = utils.Env("WINDOWS_SERVER", "iaas-module")
	kPort = utils.Env("PLAZA_PORT", "9090")
}
