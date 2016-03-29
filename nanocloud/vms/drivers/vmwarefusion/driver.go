package vmwarefusion

import (
	"os"
	"path"

	"github.com/Nanocloud/community/nanocloud/vms"
)

type driver struct{}

func (d *driver) Open(options map[string]string) (vms.VM, error) {

	err := os.MkdirAll(path.Join(options["STORAGE_DIR"], "vm"), 0755)
	if err != nil {
		return nil, err
	}

	return &vm{
		plazaLocation: options["PLAZA_LOCATION"],
		storageDir:    options["STORAGE_DIR"],
	}, nil
}
