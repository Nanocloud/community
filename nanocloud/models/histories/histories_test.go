package histories

import (
	"log"
	"testing"
	"time"

	"github.com/Nanocloud/community/nanocloud/models/users"
)

var (
	total        = 0
	history_num  = 0
	user         = &users.User{}
	connectionId = "fake-connection-id"
	startDate    []string
	endDate      []string
)

func init() {
	new_user, err := users.CreateUser(
		true,
		"new@nanocloud.com",
		"Test",
		"user",
		"secret",
		false,
	)

	if err != nil {
		log.Panicln("Can't create new account:", err.Error())
	}
	if new_user == nil {
		log.Panicln("Can't create new account")
	}
	user = new_user
}

func countEntries() {
	histories, err := FindAll()
	if err != nil {
		log.Panicln("Can't retreive histories:", err.Error())
	}
	for range histories {
		total++
	}
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

	startDate = append(startDate, time.Now().Format(time.RFC3339))
	endDate = append(endDate, time.Now().Format(time.RFC3339))
	_, err := CreateHistory(user.GetID(), user.Email, user.FirstName, user.LastName, connectionId, startDate[history_num], endDate[history_num])
	if err != nil {
		log.Panicln("Can't add historic:", err.Error())
	}

	if total-expected_num_rows+1 != 0 {
		t.Errorf("Unexpected number of rows returned: Expected %d, have %d", expected_num_rows+1, total)
	}

	err = users.DeleteUser(user.GetID())
	if err != nil {
		t.Errorf("Can't delete user: %s\n", err.Error())
	}
}
