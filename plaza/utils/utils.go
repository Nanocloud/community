// +build windows,amd64

package utils

import (
	"syscall"
	"unsafe"

	"github.com/labstack/gommon/log"
)

type startupinfo struct {
	/* DWORD */ cb uint32
	/* LPSTR */ lpReserved uintptr
	/* LPSTR */ lpDesktop uintptr
	/* LPSTR */ lpTitle uintptr
	/* DWORD */ dwX uint32
	/* DWORD */ dwY uint32
	/* DWORD */ dwXSize uint32
	/* DWORD */ dwYSize uint32
	/* DWORD */ dwXCountChars uint32
	/* DWORD */ dwYCountChars uint32
	/* DWORD */ dwFillAttribute uint32
	/* DWORD */ dwFlags uint32
	/* WORD */ wShowWindow uint16
	/* WORD */ cbReserved2 uint16
	/* LPBYTE */ lpReserved2 uintptr
	/* HANDLE */ hStdInput uintptr
	/* HANDLE */ hStdOutput uintptr
	/* HANDLE */ hStdError uintptr
}

type processinfo struct {
	/* HANDLE */ hProcess uintptr
	/* HANDLE */ hThread uintptr
	/* DWORD */ dwProcessId uint32
	/* DWORD */ dwThreadId uint32
}

type HANDLE uintptr
type PHANDLE *HANDLE

const (
	LOGON_WITH_PROFILE        = 0x1
	LOGON32_LOGON_BATCH       = 4
	LOGON32_PROVIDER_DEFAULT  = 0
	LOGON32_LOGON_INTERACTIVE = 2
)

func ExecuteCommandAsAdmin(cmd, username, pwd, domain string) {
	var si startupinfo
	var handle HANDLE
	var pi processinfo

	si.cb = uint32(unsafe.Sizeof(si))

	a := syscall.MustLoadDLL("advapi32.dll")
	LogonUserW := a.MustFindProc("LogonUserW")
	r1, r2, lastError := LogonUserW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(domain))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pwd))),
		LOGON32_LOGON_INTERACTIVE,
		LOGON32_PROVIDER_DEFAULT,
		uintptr(unsafe.Pointer(&handle)),
	)
	log.Error(r1)
	log.Error(r2)
	log.Error(lastError)

	CreateProcessAsUser := a.MustFindProc("CreateProcessAsUserW")
	r1, r2, lastError = CreateProcessAsUser.Call(
		uintptr(unsafe.Pointer(handle)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cmd))),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(0),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(&si)),
		uintptr(unsafe.Pointer(&pi)),
	)
	log.Error(r1)
	log.Error(r2)
	log.Error(lastError)

	b := syscall.MustLoadDLL("Kernel32.dll")
	CloseHandle := b.MustFindProc("CloseHandle")
	r1, r2, lastError = CloseHandle.Call(
		uintptr(unsafe.Pointer(handle)),
	)
	log.Error(r1)
	log.Error(r2)
	log.Error(lastError)
}
