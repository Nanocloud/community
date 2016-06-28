package histories

import (
	"errors"
	"github.com/Nanocloud/community/nanocloud/connectors/db"
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
		`SELECT id, userid, usermail, userfirstname,
		userlastname, connectionid, startdate, enddate FROM histories`,
	)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		h := History{}

		res.Scan(
			&h.Id,
			&h.UserId,
			&h.UserMail,
			&h.UserFirstname,
			&h.UserLastname,
			&h.ConnectionId,
			&h.StartDate,
			&h.EndDate,
		)

		result = append(result, &h)
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
