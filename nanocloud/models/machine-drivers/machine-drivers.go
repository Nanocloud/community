package machinedrivers

import (
	"github.com/Nanocloud/community/nanocloud/connectors/vms"
	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/labstack/gommon/log"
	"github.com/manyminds/api2go/jsonapi"
)

type MachineDriver struct {
	ID string
}

func (d *MachineDriver) GetID() string {
	return d.ID
}

func (d *MachineDriver) SetID(id string) error {
	d.ID = id
	return nil
}

func (d *MachineDriver) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "machine-types",
			Name: "types",
		},
	}
}

func (d *MachineDriver) GetReferencedIDs() []jsonapi.ReferenceID {
	types, err := vms.Types()
	if err != nil {
		log.Error(err)
		return nil
	}

	rt := make([]jsonapi.ReferenceID, 0)

	for _, t := range types {

		rt = append(rt, jsonapi.ReferenceID{
			ID:   t.GetID(),
			Type: "machine-types",
			Name: "types",
		})
	}

	return rt
}

func (d *MachineDriver) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	types, err := vms.Types()
	if err != nil {
		log.Error(err)
		return nil
	}

	rt := make([]jsonapi.MarshalIdentifier, 0)

	for _, t := range types {
		rt = append(rt, t)
	}

	return rt
}

func FindAll() ([]*MachineDriver, error) {
	drivers := make([]*MachineDriver, 1)
	drivers[0] = &MachineDriver{
		ID: utils.Env("IAAS", ""),
	}
	return drivers, nil
}
