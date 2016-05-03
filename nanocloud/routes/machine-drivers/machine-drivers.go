package machinedrivers

import (
	"net/http"

	machinedrivers "github.com/Nanocloud/community/nanocloud/models/machine-drivers"
	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/labstack/echo"
)

func FindAll(c *echo.Context) error {
	drivers, err := machinedrivers.FindAll()
	if err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, drivers)
}
