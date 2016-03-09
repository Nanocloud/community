package qemu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type vm struct {
}

func (v *vm) Types() ([]vms.MachineType, error) {
	return []vms.MachineType{defaultType}, nil
}

func (v *vm) Create(name, password string, t vms.MachineType) (vms.Machine, error) {

	m := machine{id: name}
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/iaas/"+m.Id()+"/download", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return nil, err
	}
	return &m, nil
}

func (v *vm) Machines() ([]vms.Machine, error) {
	iaas := os.Getenv("WIN_SERVER")
	resp, err := http.Get("http://" + iaas + ":8080/api/iaas")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	type Status struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes VmInfo
	}
	var State struct {
		Data []Status `json:"data"`
	}

	err = json.Unmarshal(body, &State)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var machines = make([]vms.Machine, len(State.Data))
	for i, val := range State.Data {
		machines[i] = &machine{id: val.Id}
	}
	return machines, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	return &machine{id: id}, nil
}
