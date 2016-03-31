package router

import (
	"net/http"

	"github.com/Nanocloud/community/plaza/routes/about"
	"github.com/Nanocloud/community/plaza/routes/apps"
	"github.com/Nanocloud/community/plaza/routes/files"
	"github.com/Nanocloud/community/plaza/routes/power"
	"github.com/Nanocloud/community/plaza/routes/sessions"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func Start() {
	e := echo.New()

	e.Get("/", about.Get)

	/***
	FILES
	***/

	e.Get("/files", files.Get)
	e.Post("/upload", files.Post)
	e.Get("/upload", files.GetUpload)

	/***
	POWER
	***/

	e.Get("/shutdown", power.ShutDown)
	e.Get("/restart", power.Restart)
	e.Get("/checkrds", power.CheckRDS)

	/***
	SESSIONS
	***/

	e.Get("/sessions/:id", sessions.Get)
	e.Delete("/sessions/:id", sessions.Logoff)

	/***
	APPS
	***/

	e.Post("/publishapp", apps.PublishApp)
	e.Get("/apps", apps.GetApps)
	e.Delete("/apps/:id", apps.UnpublishApp)

	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		c.JSON(
			http.StatusInternalServerError,
			hash{
				"error": err.Error(),
			},
		)
	})

	log.Println("Listenning on port: 9090")
	e.Run(":9090")
}
