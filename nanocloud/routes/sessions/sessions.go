package sessions

import (
	"io/ioutil"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

var kServer string

type hash map[string]interface{}

func List(c *echo.Context) error {
	user := c.Get("user").(*users.User)

	resp, err := http.Get("http://" + kServer + ":" + utils.Env("PLAZA_PORT", "9090") + "/sessions/" + user.Sam)
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

func Logoff(c *echo.Context) error {
	user := c.Get("user").(*users.User)

	req, err := http.NewRequest("DELETE", "http://"+kServer+":"+utils.Env("PLAZA_PORT", "9090")+"/sessions/"+user.Sam, nil)
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
	kServer = utils.Env("SERVER", "localhost")
}
