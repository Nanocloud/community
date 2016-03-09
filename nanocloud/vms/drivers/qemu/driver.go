package qemu

import "github.com/Nanocloud/community/nanocloud/vms"

type driver struct{}

func (d *driver) Open(options map[string]string) (vms.VM, error) {
	return &vm{}, nil
}
