package logout

import (
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"gopkg.in/labstack/echo.v1"
)

func Post(c *echo.Context) error {
	user := c.Get("user").(*users.User)
	req := c.Request()
	accessToken, _ := oauth2.GetAuthorizationHeaderValue(req, "Bearer")

	_, err := db.Exec(
		`DELETE FROM oauth_access_tokens
		WHERE token = $1::varchar
		AND user_id = $2::varchar`,
		accessToken,
		user.Id,
	)
	return err
}
