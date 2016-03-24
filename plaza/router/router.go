package router

import (
	"net/http"

	"github.com/Nanocloud/community/plaza/routes/about"
	"github.com/Nanocloud/community/plaza/routes/apps"
	"github.com/Nanocloud/community/plaza/routes/files"
	"github.com/Nanocloud/community/plaza/routes/power"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func Start() {
	e := echo.New()

	e.Get("/", about.Get)
	e.Get("/files", files.Get)
	e.Post("/upload", files.Post)
	e.Get("/upload", files.GetUpload)

	e.Get("/shutdown", power.ShutDown)
	e.Get("/restart", power.Restart)
	e.Get("/checkrds", power.CheckRDS)

	/***
	APPS
	***/

	e.Post("/publishapp", apps.PublishApp)
	e.Get("/apps", apps.GetApps)

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
