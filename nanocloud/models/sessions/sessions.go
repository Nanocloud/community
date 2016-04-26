package sessions

import (
	"encoding/json"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var kServer string

type hash map[string]interface{}

func GetAll(userSam string) ([]Session, error) {

	var sessionList []Session

	resp, err := http.Get("http://" + kServer + ":" + utils.Env("PLAZA_PORT", "9090") + "/sessions/" + userSam)

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

		var session Session
		for t, s := range tab {
			log.Info(s)
			if t == 0 {
				session.SessionName = s
			}
			if t == 1 {
				session.Username = s
			}
			if t == 2 {
				session.Id = s
			}
			if t == 3 {
				session.State = s
			}
		}
		sessionList = append(sessionList, session)
	}
	return sessionList, nil
}

func init() {
	kServer = utils.Env("SERVER", "localhost")
}
