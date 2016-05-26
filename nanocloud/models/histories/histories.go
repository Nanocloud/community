package histories

import (
	"errors"
	"strings"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/apps"
	"github.com/Nanocloud/community/nanocloud/models/users"
	uuid "github.com/satori/go.uuid"
)

var (
	HistoryNotCreated = errors.New("history not created")
)

func isAlphaNum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isValidId(id string) bool {
	l := len(id)

	for i := 0; i < l; i++ {
		c := id[i]
		if c != '-' && !isAlphaNum(c) {
			return false
		}
	}
	return true
}

func escapeId(id string) string {
	return "'" + id + "'"
}

func FindAll() ([]*History, error) {

	result := make([]*History, 0)

	res, err := db.Query(
		`SELECT histories.id, histories.userid, histories.usermail, histories.userfirstname,
		histories.userlastname, apps.id as app_id, histories.startdate, histories.enddate
		FROM histories
		LEFT JOIN apps ON apps.alias = histories.connectionid`,
	)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	userList := make(map[string][]*History)
	userIds := make([]string, 0)

	appList := make(map[string][]*History)
	appIds := make([]string, 0)

	for res.Next() {
		h := History{}

		var userId string
		var appId string

		res.Scan(
			&h.Id,
			&userId,
			&h.UserMail,
			&h.UserFirstname,
			&h.UserLastname,
			&appId,
			&h.StartDate,
			&h.EndDate,
		)

		userHistories := userList[userId]
		if userHistories == nil {
			userHistories = make([]*History, 0)
		}
		if isValidId(userId) {
			userList[userId] = append(userHistories, &h)
			userIds = append(userIds, escapeId(userId))
		}

		appHistories := appList[appId]
		if appHistories == nil {
			appHistories = make([]*History, 0)
		}
		if isValidId(appId) {
			appList[appId] = append(appHistories, &h)
			appIds = append(appIds, escapeId(appId))
		}
		result = append(result, &h)
	}

	if len(userIds) > 0 {
		sUserIds := strings.Join(userIds, ",")

		rows, err := db.Query(
			`SELECT id, activated,
			email,
			first_name, last_name,
			is_admin
			FROM users
			WHERE id in (` + sUserIds + ")",
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			user := users.User{}
			rows.Scan(
				&user.Id,
				&user.Activated,
				&user.Email,
				&user.FirstName,
				&user.LastName,
				&user.IsAdmin,
			)

			for _, h := range userList[user.Id] {
				h.user = &user
			}
		}
	}

	if len(appIds) > 0 {

		sAppIds := strings.Join(appIds, ",")

		rows, err := db.Query(
			`SELECT id, collection_name,
			alias, display_name,
			file_path,
			icon_content
			FROM apps
			WHERE id in (` + sAppIds + ")",
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			app := apps.Application{}
			rows.Scan(
				&app.Id,
				&app.CollectionName,
				&app.Alias,
				&app.DisplayName,
				&app.FilePath,
				&app.IconContents,
			)
			for _, h := range appList[app.Id] {
				h.application = &app
			}
		}
	}

	return result, nil
}

func CreateHistory(
	userId string,
	userMail string,
	userFirstname string,
	userLastname string,
	connectionId string,
	startDate string,
	endDate string,
) (*History, error) {
	id := uuid.NewV4().String()

	rows, err := db.Query(
		`INSERT INTO histories
		(id, userid, usermail, userfirstname, userlastname, connectionid, startdate, enddate)
		VALUES(	$1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::varchar, $7::varchar, $8::varchar)`,
		id, userId, userMail, userFirstname, userLastname, connectionId, startDate, endDate)

	if err != nil {
		return nil, err
	}

	rows.Close()

	rows, err = db.Query(
		`SELECT id, userid, usermail, userfirstname, userlastname, connectionid, startdate, enddate
		FROM histories WHERE id = $1::varchar`, id)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, HistoryNotCreated
	}

	var history History
	rows.Scan(
		&history.Id,
		&history.UserId,
		&history.UserMail,
		&history.UserFirstname,
		&history.UserLastname,
		&history.ConnectionId,
		&history.StartDate,
		&history.EndDate,
	)

	rows.Close()

	return &history, err
}
