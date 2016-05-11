// +build windows

package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	logon32ProviderDefault   = 0
	logonWithProfile         = 1
	logon32LogonInteractive  = 2
	logon32LogonBatch        = 4
	uintptrFlagInherit       = 0x00000001
	createUnicodeEnvironment = 0x00000400

	tokenQuery            = 0x0008
	tokenAdjustPrivileges = 0x0020

	handleFlagInherit = 0x00000001

	startfUseStdHandles = 0x00000100
)

func impersonateLoggedOnUser(token syscall.Handle) error {
	proc, err := loadProc("Advapi32.dll", "ImpersonateLoggedOnUser")
	if err != nil {
		return err
	}

	r1, _, err := proc.Call(uintptr(token))
	if r1 == 0 {
		fmt.Println(err)
		return err
	}
	return nil
}

func getUserProfileDirectory(token syscall.Handle) (*uint16, error) {
	proc, err := loadProc("Userenv.dll", "GetUserProfileDirectoryW")
	if err != nil {
		return nil, err
	}

	buffSize := (260 + 1) // (MAX_PATH) * sizeof(WCHAR)
	buff := make([]uint16, buffSize)
	r1, _, err := proc.Call(
		uintptr(token),
		uintptr(unsafe.Pointer(&buff[0])),
		uintptr(unsafe.Pointer(&buffSize)),
	)
	if r1 != 1 {
		return nil, err
	}
	return &buff[0], nil
}

func openProcessToken(handle uintptr, desiredAccess uint32, tokenHandle *uintptr) error {
	proc, err := loadProc("Advapi32.dll", "OpenProcessToken")
	if err != nil {
		return err
	}

	r1, _, err := proc.Call(
		handle,
		uintptr(desiredAccess),
		uintptr(unsafe.Pointer(tokenHandle)),
	)
	if r1 != 1 {
		return err
	}
	return nil
}

func destroyEnvironmentBlock(env *uint16) error {
	proc, err := loadProc("Userenv.dll", "DestroyEnvironmentBlock")
	if err != nil {
		return err
	}
	r1, _, err := proc.Call(uintptr(unsafe.Pointer(env)))
	if r1 == 0 {
		return err
	}
	return nil
}

func createEnvironmentBlock(token syscall.Handle, inherit bool) ([]uint16, error) {
	proc, err := loadProc("Userenv.dll", "CreateEnvironmentBlock")
	if err != nil {
		return nil, err
	}

	iInherit := 0
	if inherit {
		iInherit = 1
	}

	var env *uint16

	r1, _, err := proc.Call(
		uintptr(unsafe.Pointer(env)),
		uintptr(token),
		uintptr(iInherit),
	)

	if r1 == 1 {
		l := 0
		for l = 0; *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(env)) + uintptr(l))) != 0; l++ {
		}
		rt := make([]uint16, l)

		for i := 0; i < l; i++ {
			rt[i] = *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(env)) + uintptr(i)))
		}

		err = destroyEnvironmentBlock(env)
		if err != nil {
			return nil, err
		}

		return rt, nil
	}
	return nil, err
}

func createProcessWithLogon(
	username string,
	domain string,
	password string,
	logonFlags uint32,
	applicationName string,
	cmd string,
	creationFlags uint32,
	environment uintptr,
	currentDirectory string,
	si *syscall.StartupInfo,
	pi *syscall.ProcessInformation,
) error {
	fmt.Println("createProcessWithLogon")
	fmt.Println(cmd)
	dll, err := syscall.LoadDLL("advapi32.dll")
	if err != nil {
		return err
	}
	proc, err := dll.FindProc("CreateProcessWithLogonW")
	if err != nil {
		return err
	}

	r1, _, err := proc.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(domain))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(password))),
		uintptr(logonFlags),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(applicationName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cmd))),
		uintptr(creationFlags),
		environment,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(currentDirectory))),
		uintptr(unsafe.Pointer(si)),
		uintptr(unsafe.Pointer(pi)),
	)
	if r1 == 1 {
		return nil
	}
	return err
}

func logonUser(username, domain, password string, logonType, logonProvider uint32) (hd syscall.Handle, err error) {
	dll, err := loadDLL("advapi32.dll")
	if err != nil {
		return
	}
	proc, err := dll.FindProc("LogonUserW")
	if err != nil {
		return
	}
	r1, _, err := proc.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(domain))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(password))),
		uintptr(logonType),
		uintptr(logonProvider),
		uintptr(unsafe.Pointer(&hd)),
	)
	if r1 == 1 {
		err = nil
	}
	return
}

func createProcessAsUser(
	token syscall.Handle,
	applicationName *uint16,
	cmd *uint16,
	procSecurity *syscall.SecurityAttributes,
	threadSecurity *syscall.SecurityAttributes,
	inheritHandles bool,
	creationFlags uint32,
	environment *uint16,
	currentDir *uint16,
	startupInfo *syscall.StartupInfo,
	outProcInfo *syscall.ProcessInformation,
) error {
	dll, err := loadDLL("advapi32.dll")
	if err != nil {
		return err
	}
	proc, err := dll.FindProc("CreateProcessAsUserW")
	if err != nil {
		return err
	}

	iInheritHandles := 0
	if inheritHandles {
		iInheritHandles = 1
	}

	r1, _, err := proc.Call(
		uintptr(token),
		uintptr(unsafe.Pointer(applicationName)),
		uintptr(unsafe.Pointer(cmd)),
		uintptr(unsafe.Pointer(procSecurity)),
		uintptr(unsafe.Pointer(threadSecurity)),
		uintptr(iInheritHandles),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(environment)),
		uintptr(unsafe.Pointer(currentDir)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(outProcInfo)),
	)

	if r1 == 1 {
		return nil
	}
	return err
}
