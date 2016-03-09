package vms

type Driver interface {
	Open(options map[string]string) (VM, error)
}
