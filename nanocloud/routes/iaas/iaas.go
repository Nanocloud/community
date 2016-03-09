package iaas

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/connectors/vms"
	vm "github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

type JsonMachine struct {
	Id     string
	Name   string
	Status string
	Ip     string
}

func MachinetoStruct(rawmachine vm.Machine) JsonMachine {
	var mach JsonMachine
	mach.Id = rawmachine.Id()
	mach.Name, _ = rawmachine.Name()
	status, _ := rawmachine.Status()
	mach.Status = vm.StatusToString(status)
	ip, _ := rawmachine.IP()
	mach.Ip = string(ip)
	return mach
}

func ListRunningVM(c *echo.Context) error {
	machines, err := vms.Machines()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, hash{
				"errors": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
	}
	type Attr struct {
		Name   string `json:"name"`
		Ip     string `json:"ip"`
		Status string `json:"status"`
		Id     string `json:"id"`
	}
	type VM struct {
		Id  string `json:"id"`
		Att Attr   `json:"attributes"`
	}
	var res = make([]VM, len(machines))
	for i, val := range machines {
		res[i].Att.Name, err = val.Name()
		if err != nil {
			log.Error(err)
		}
		res[i].Att.Id = val.Id()
		status, err := val.Status()
		if err != nil {
			log.Error(err)
		}
		res[i].Att.Status = vm.StatusToString(status)
		ip, _ := val.IP()
		res[i].Att.Ip = string(ip)
		if err != nil {
			log.Error(err)
		}
	}

	return c.JSON(http.StatusOK, hash{"data": res})
}

func StopVM(c *echo.Context) error {
	machine, err := vms.Machine(c.Param("id"))

	err = machine.Stop()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, hash{
				"errors": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": MachinetoStruct(machine),
		})
}

func StartVM(c *echo.Context) error {
	machine, err := vms.Machine(c.Param("id"))

	err = machine.Start()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, hash{
				"errors": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": MachinetoStruct(machine),
		})
}

func CreateVM(c *echo.Context) error {
	//TODO READ BODY TO GET PASSWORD AND TYPE
	vm, err := vms.Create(c.Param("id"), "", nil)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, hash{
				"errors": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": MachinetoStruct(vm),
		})
}
