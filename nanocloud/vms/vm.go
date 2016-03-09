package vms

type VM interface {
	Machines() ([]Machine, error)
	Machine(id string) (Machine, error)
	Create(name, password string, t MachineType) (Machine, error)
	Types() ([]MachineType, error)
}
