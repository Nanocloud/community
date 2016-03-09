package qemu

import (
	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type driver struct{}

func (d *driver) Open(options map[string]string) (vms.VM, error) {
	log.Error("SECONDARY OPEN")
	return &vm{}, nil
}
