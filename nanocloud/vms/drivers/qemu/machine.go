package qemu

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type machine struct {
	id string
}

type VmInfo struct {
	Ico         string `json:"ico"`
	Name        string `json:"-"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Locked      bool   `json:"locked"`
	CurrentSize string `json:"current_size"`
	TotalSize   string `json:"total_size"`
}

func (m *machine) Status() (vms.MachineStatus, error) {
	ip, _ := m.IP()
	resp, err := http.Get("http://" + string(ip) + ":8080/api/iaas")
	if err != nil {
		log.Error(err)
		return vms.StatusUnknown, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return vms.StatusUnknown, err
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
		return vms.StatusUnknown, err
	}
	for _, val := range State.Data {
		switch val.Attributes.Status {
		case "running":
			return vms.StatusUp, nil
		case "booting":
			return vms.StatusBooting, nil
		case "download":
			return vms.StatusDownloading, nil
		case "available":
			return vms.StatusDown, nil
		}
		return vms.StatusUnknown, nil
	}
	return vms.StatusUnknown, nil
}

func (m *machine) IP() (net.IP, error) {
	iaas := os.Getenv("WIN_SERVER")
	return []byte(iaas), nil
}

func (m *machine) Type() (vms.MachineType, error) {
	return defaultType, nil
}

func (m *machine) Start() error {
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/iaas/"+m.id+"/start", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return err
	}
	return nil
}

func (m *machine) Stop() error {
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/iaas/"+m.id+"/stop", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return err
	}
	return nil
}

func (m *machine) Terminate() error {
	return nil
}

func (m *machine) Id() string {
	return m.id
}

func (m *machine) Name() (string, error) {
	return "Windows Active Directory", nil
}
