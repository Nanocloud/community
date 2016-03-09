package vms

import "net"

type MachineStatus int

const (
	StatusUnknown     MachineStatus = 0
	StatusDown        MachineStatus = 1
	StatusUp          MachineStatus = 2
	StatusTerminated  MachineStatus = 3
	StatusBooting     MachineStatus = 4
	StatusDownloading MachineStatus = 5
)

type Machine interface {
	Id() string
	Name() (string, error)
	Status() (MachineStatus, error)
	IP() (net.IP, error)
	Type() (MachineType, error)

	Start() error
	Stop() error
	Terminate() error
}

func StatusToString(status MachineStatus) string {
	switch status {
	case StatusDown:
		return "available"
	case StatusUp:
		return "running"
	case StatusTerminated:
		return "terminated"
	case StatusBooting:
		return "booting"
	case StatusDownloading:
		return "download"
	}
	return "unknown"
}
