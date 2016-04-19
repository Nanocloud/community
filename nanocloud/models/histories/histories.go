package histories

import (
	"errors"
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	uuid "github.com/satori/go.uuid"
)

var (
	HistoryNotCreated = errors.New("history not created")
)

func GetAll() ([]*History, error) {

	var historyList []*History
	rows, err := db.Query(
		`SELECT id, userid, connectionid,
		startdate, enddate
		FROM histories`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		history := History{}

		rows.Scan(
			&history.Id,
			&history.UserId,
			&history.ConnectionId,
			&history.StartDate,
			&history.EndDate,
		)
		historyList = append(historyList, &history)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	rows.Close()

	return historyList, nil
}

func CreateHistory(
	userId string,
	connectionId string,
	startDate string,
	endDate string,
) (*History, error) {

	id := uuid.NewV4().String()

	rows, err := db.Query(
		`INSERT INTO histories
		(id, userid, connectionid, startdate, enddate)
		VALUES(	$1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar)`,
		id, userId, connectionId, startDate, endDate)

	if err != nil {
		return nil, err
	}

	rows.Close()

	rows, err = db.Query(
		`SELECT id, userid, connectionid, startdate, enddate
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
		&history.ConnectionId,
		&history.StartDate,
		&history.EndDate,
	)

	rows.Close()

	return &history, err
}
