package qemu

import "github.com/Nanocloud/community/nanocloud/vms"

const (
	windowsIsoUri = "http://care.dlservice.microsoft.com/dl/download/6/2/A/62A76ABB-9990-4EFC-A4FE-C7D698DAEB96/9600.17050.WINBLUE_REFRESH.140317-1640_X64FRE_SERVER_EVAL_EN-US-IR3_SSS_X64FREE_EN-US_DV9.ISO"
)

func init() {
	vms.Register("qemu", &driver{})
}
