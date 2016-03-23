package router

import (
	"net/http"

	"github.com/Nanocloud/community/plaza/server/routes/about"
	"github.com/Nanocloud/community/plaza/server/routes/apps"
	"github.com/Nanocloud/community/plaza/server/routes/files"
	"github.com/Nanocloud/community/plaza/server/routes/power"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func Start() {
	e := echo.New()

	e.Get("/", about.Get)
	e.Get("/files", files.Get)
	e.Post("/upload", files.Post)

	e.Get("/shutdown", power.ShutDown)
	e.Get("/restart", power.Restart)
	e.Get("/checkrds", power.CheckRDS)

	/***
	APPS
	***/

	e.Post("/publishapp", apps.PublishApp)
	e.Get("/apps", apps.GetApps)

	/***
	PROVISIONING
	***/
	/*
		e.Post("/disablewu", prov.DisableWU)
		e.Get("/disablewu", prov.CheckWU)
		e.Post("/installad", prov.InstallAD)
		e.Get("/installad", prov.CheckAD)
		e.Post("/enablerdp", prov.EnableRDP)
		e.Get("/enablerdp", prov.CheckRDP)
		e.Post("/installrds", prov.InstallRDS)
		e.Get("/installrds", prov.CheckRDS)
		e.Post("/createou", prov.CreateOU)
		e.Get("/createou", prov.CheckOU)
		e.Post("/installadcs", prov.InstallADCS)
		e.Get("/installadcs", prov.CheckADCS)
		e.Post("/sessiondeploy", prov.SessionDeploy)
		e.Get("/sessiondeploy", prov.CheckCollection)*/

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
