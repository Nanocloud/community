package qemu

import "github.com/Nanocloud/community/nanocloud/vms"

func init() {
	vms.Register("qemu", &driver{})
}
