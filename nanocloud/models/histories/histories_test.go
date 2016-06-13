package histories

import (
	"log"
	"testing"
	"time"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
)

var (
	history_num  = 0
	user         = &users.User{}
	connectionId = "fake-connection-id"
	startDate    []string
	endDate      []string
)

func init() {
	admin_user, err := users.GetUserFromEmailPassword("admin@nanocloud.com", "Nanocloud123+")
	if err != nil {
		log.Panicf("Can't retreive administrator account: %s\r\n", err.Error())
	}
	if admin_user == nil {
		log.Panicf("Can't retreive administrator account\r\n")
	}
	user = admin_user
}

func TestCreateHistory(t *testing.T) {
	startDate = append(startDate, time.Now().Format(time.RFC3339))
	endDate = append(endDate, time.Now().Format(time.RFC3339))
	history, err := CreateHistory(user.GetID(), user.Email, user.FirstName, user.LastName, connectionId, startDate[history_num], endDate[history_num])
	if err != nil {
		t.Errorf("Cannot create history: %s", err.Error())
	}

	switch {
	case history.Id == "":
		t.Errorf("'history.Id' field doesn't match the inserted value")
	case history.UserId != user.GetID():
		t.Errorf("'history.UserId' field doesn't match the inserted value")
	case history.UserMail != user.Email:
		t.Errorf("'user.Email' field doesn't match the inserted value")
	case history.UserFirstname != user.FirstName:
		t.Errorf("'user.FirstName' field doesn't match the inserted value")
	case history.UserLastname != user.LastName:
		t.Errorf("'user.LastName' field doesn't match the inserted value")
	case history.ConnectionId != connectionId:
		t.Errorf("'history.ConnectionId' field doesn't match the inserted value")
	case history.StartDate != startDate[history_num]:
		t.Errorf("'history.StartDate' field doesn't match the inserted value")
	case history.EndDate != endDate[history_num]:
		t.Errorf("'history.EndDate' field doesn't match the inserted value")
	}
	history_num++
}

func TestFindAll(t *testing.T) {
	var expected_num_rows int = history_num
	var num_rows int

	startDate = append(startDate, time.Now().Format(time.RFC3339))
	endDate = append(endDate, time.Now().Format(time.RFC3339))
	_, err := CreateHistory(user.GetID(), user.Email, user.FirstName, user.LastName, connectionId, startDate[history_num], endDate[history_num])
	if err != nil {
		log.Panicf("Can't add historic: %s", err.Error())
	}
	rows, err := db.Query("SELECT COUNT(*) FROM histories")
	if err != nil {
		t.Errorf("Can't count histories")
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&num_rows)
		if err != nil {
			t.Errorf("Error when trying to scan query result: %s", err.Error())
		}
		if num_rows != expected_num_rows+1 {
			t.Errorf("Unexpected number of rows returned: Expected %d, have %d", expected_num_rows+1, num_rows)
		}
	} else {
		t.Errorf("No result was returned by the query")
	}
	for _, row := range startDate {
		db.Exec(`DELETE FROM histories where startDate=$1::varchar`, row)
	}
}
