// +build windows

package windows

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func LaunchApp(name []string) (uint32, error) {
	sessions, err := wtsEnumerateSessions(0)
	if err != nil {
		return 0, err
	}

	var session *wtsSessionInfo

	for _, s := range sessions {
		if strings.Index(s.winStationName, "RDP-Tcp#") == 0 && s.state == wtsActive {
			session = &s
			break
		}
	}

	if session == nil {
		return 0, errors.New("no active session found")
	}

	token, err := wtsQueryUserToken(session.sessionID)
	if err != nil {
		return 0, fmt.Errorf("Query User Token Failed: %s", err.Error())
	}
	defer syscall.CloseHandle(token)

	cmd := strings.Join(name, "")

	id := strconv.FormatInt(int64(session.sessionID), 10)

	wsName, err := syscall.UTF16PtrFromString(`C:\PSTools\PsExec.exe -d -i ` + id + " " + cmd)
	if err != nil {
		return 0, err
	}

	NORMAL_PRIORITY_CLASS := 0x00000020
	CREATE_NEW_CONSOLE := 0x00000010

	var flags uint32
	flags |= syscall.CREATE_UNICODE_ENVIRONMENT
	flags |= uint32(NORMAL_PRIORITY_CLASS)
	flags |= uint32(CREATE_NEW_CONSOLE)

	si := new(syscall.StartupInfo)
	si.Cb = uint32(unsafe.Sizeof(*si))
	wsDesktop, err := syscall.UTF16PtrFromString(`winsta0\default`)
	if err != nil {
		return 0, err
	}
	si.Desktop = wsDesktop
	pi := new(syscall.ProcessInformation)

	err = createProcessAsUser(
		token,
		nil,
		wsName,
		nil,
		nil,
		false,
		flags,
		nil,
		nil,
		si,
		pi,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateProcessAsUser Failed: %s", err.Error())
	}
	return pi.ProcessId, nil
}
