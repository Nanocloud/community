package sessions

import (
	"encoding/json"
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var kServer string
var kPort string

type hash map[string]interface{}

func GetAll(userSam string) ([]Session, error) {

	var sessionList []Session

	resp, err := http.Get("http://" + kServer + ":" + kPort + "/sessions/" + userSam)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var body struct {
		Data [][]string
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, err
	}

	for _, tab := range body.Data {

		rows, err := db.Query(
			`SELECT users.id FROM users
			left join users_windows_user on users.id = users_windows_user.user_id
			left join windows_users on users_windows_user.windows_user_id = windows_users.id
			WHERE windows_users.sam = $1::varchar`,
			tab[1])

		if err != nil {
			return nil, err
		}

		defer rows.Close()
		var user_id string
		if rows.Next() {
			err = rows.Scan(
				&user_id,
			)

			if err != nil {
				return nil, err
			}

			var session Session
			session.SessionName = tab[0]
			session.Username = tab[1]
			session.Id = tab[2]
			session.State = tab[3]
			session.UserId = user_id
			sessionList = append(sessionList, session)
		}
	}
	return sessionList, nil
}

func init() {
	kServer = utils.Env("EXECUTION_SERVERS", "iaas-module")
	kPort = utils.Env("PLAZA_PORT", "9090")
}
