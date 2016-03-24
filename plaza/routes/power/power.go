package power

import (
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
)

type hash map[string]interface{}

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

func ShutDown(c *echo.Context) error {
	cmd := exec.Command("powershell.exe", "Stop-Computer -Force")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return reterr(err, string(resp), c)
	}
	return retok(c)
}

func Restart(c *echo.Context) error {
	cmd := exec.Command("powershell.exe", "Restart-Computer -Force")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			hash{
				"error":    err.Error(),
				"response": resp,
			},
		)
	}
	return c.JSON(
		http.StatusOK,
		hash{
			"success": true,
		},
	)
}

func CheckRDS(c *echo.Context) error {
	cmd := exec.Command("powershell.exe", "Write-Host (Get-Service -Name RDMS).status")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			hash{
				"error":    err.Error(),
				"response": resp,
			},
		)
	}
	return c.JSON(
		http.StatusOK,
		hash{
			"state": string(resp),
		},
	)
}
