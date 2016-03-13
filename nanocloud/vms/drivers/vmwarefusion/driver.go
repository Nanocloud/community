package vmwarefusion

import "github.com/Nanocloud/community/nanocloud/vms"

type driver struct{}

func (d *driver) Open(options map[string]string) (vms.VM, error) {
	return &vm{
		plazaLocation: options["PLAZA_LOCATION"],
		storageDir:    options["STORAGE_DIR"],
	}, nil
}
