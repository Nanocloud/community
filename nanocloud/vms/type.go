package vms

type MachineType interface {
	Id() string
	Label() string
}
