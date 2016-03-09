package vms

import (
	"os"

	"github.com/Nanocloud/community/nanocloud/vms"
	_ "github.com/Nanocloud/community/nanocloud/vms/drivers/qemu"
	log "github.com/Sirupsen/logrus"
)

var _vm *vms.VM

func getInstance() (*vms.VM, error) {
	if _vm == nil {

		iaas := os.Getenv("IAAS")
		if len(iaas) == 0 {
			log.Fatal("No iaas provided")
		}
		var err error
		_vm, err = vms.Open(iaas, nil)
		return _vm, err
	}
	return _vm, nil
}

func Machines() ([]vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return (*vm).Machines()
}

func Machine(id string) (vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		return nil, err
	}
	return (*vm).Machine(id)
}

func Create(name, password string, t vms.MachineType) (vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		return nil, err
	}
	return (*vm).Create(name, password, t)
}
