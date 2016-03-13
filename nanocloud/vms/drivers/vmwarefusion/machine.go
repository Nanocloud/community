package vmwarefusion

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type machine struct {
	id         string
	storageDir string
}

func (m *machine) vmxPath() string {
	return path.Join(m.storageDir, "vm", m.id, "conf.vmx")
}

func trim(str string) string {
	return strings.Trim(strings.Trim(str, " "), "\t")
}

func (m *machine) parseVMX() (map[string]string, error) {
	vmxfh, err := os.Open(m.vmxPath())
	if err != nil {
		return nil, err
	}
	defer vmxfh.Close()

	vmxcontent, err := ioutil.ReadAll(vmxfh)
	if err != nil {
		return nil, err
	}

	rt := make(map[string]string, 0)
	for _, line := range strings.Split(string(vmxcontent), "\n") {
		splt := strings.SplitN(line, "=", 2)
		if len(splt) == 2 {
			key := trim(splt[0])
			value := strings.Trim(trim(splt[1]), "\"")
			rt[key] = value
		}
	}

	return rt, nil
}

func (m *machine) vmx(key string) (string, error) {
	vmx, err := m.parseVMX()
	if err != nil {
		return "", err
	}
	return vmx[key], nil
}

func (m *machine) Id() string {
	return m.id
}

func (m *machine) Name() (string, error) {
	return m.vmx("displayName")
}

func (m *machine) Status() (vms.MachineStatus, error) {
	cmd := exec.Command(vmrun, "list")
	stdout, err := cmd.Output()
	if err != nil {
		return vms.StatusUnknown, err
	}
	out := string(stdout)

	splt := strings.Split(out, "\n")
	if len(splt) < 1 {
		return vms.StatusUnknown, errors.New("Unable to list running VM")
	}
	splt = splt[1:]

	vmxPath := m.vmxPath()
	for _, line := range splt {
		if line == vmxPath {
			return vms.StatusUp, nil
		}
	}
	_, err = os.Stat(vmxPath)
	if err != nil {
		if os.IsNotExist(err) {
			return vms.StatusTerminated, nil
		}
		return vms.StatusUnknown, err
	}

	return vms.StatusDown, nil
}

func (m *machine) IP() (net.IP, error) {
	// Look for generatedAddress as we're passing a VMX with addressType = "generated".
	mac, err := m.vmx("ethernet0.generatedAddress")
	if err != nil {
		return nil, err
	}

	if len(mac) > 0 {
		ip, err := getIpFromMAC(strings.ToLower(mac))
		if err != nil {
			return nil, err
		}
		if len(ip) > 0 {
			return net.ParseIP(ip), nil
		}
	}

	log.WithFields(log.Fields{
		"VM": m.id,
	}).Error("couldn't find MAC address in VMX file")

	return nil, nil
}

func (m *machine) Type() (vms.MachineType, error) {
	return defaultType, nil
}

func (m *machine) Start() error {
	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("Starting VM")

	vmx := path.Join(m.storageDir, "vm", m.id, "conf.vmx")

	log.Debugln("Executing:", vmrun, "start", vmx)
	cmd := exec.Command(vmrun, "start", vmx)
	err := cmd.Start()

	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("VM Started")

	return err
}

func (m *machine) Stop() error {
	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("Stopping VM")

	vmx := path.Join(m.storageDir, "vm", m.id, "conf.vmx")

	log.Debugln("Executing:", vmrun, "stop", vmx)
	cmd := exec.Command(vmrun, "stop", vmx)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithFields(log.Fields{
			"VM": m.id,
		}).Error(string(out))

		return err
	}

	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("VM Stopped")
	return nil
}

func (m *machine) Terminate() error {
	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("Deleting VM")

	vmx := path.Join(m.storageDir, "vm", m.id, "conf.vmx")

	log.Debugln("Executing:", vmrun, "deleteVM", vmx)
	cmd := exec.Command(vmrun, "deleteVM", vmx)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithFields(log.Fields{
			"VM": m.id,
		}).Error(string(out))
		return err
	}

	err = os.RemoveAll(path.Join(m.storageDir, "vm", m.id))

	if err != nil {
		log.WithFields(log.Fields{
			"VM": m.id,
		}).Error(err)
		return err
	}

	log.WithFields(log.Fields{
		"VM": m.id,
	}).Info("VM Deleted")

	return nil
}

func getIpFromMAC(macaddr string) (string, error) {
	// DHCP lease table for NAT vmnet interface
	leasesFiles, _ := filepath.Glob("/var/db/vmware/*.leases")
	for _, dhcpfile := range leasesFiles {
		log.Debugf("Trying to find IP address in leases file: %s", dhcpfile)
		if ipaddr, err := getIPfromDHCPLeaseFile(dhcpfile, macaddr); err == nil {
			return ipaddr, err
		}
	}

	return "", fmt.Errorf("IP not found for MAC %s in DHCP leases", macaddr)
}

func getIPfromDHCPLeaseFile(dhcpfile, macaddr string) (string, error) {
	var dhcpfh *os.File
	var dhcpcontent []byte
	var lastipmatch string
	var currentip string
	var lastleaseendtime time.Time
	var currentleadeendtime time.Time
	var err error

	if dhcpfh, err = os.Open(dhcpfile); err != nil {
		return "", err
	}
	defer dhcpfh.Close()

	if dhcpcontent, err = ioutil.ReadAll(dhcpfh); err != nil {
		return "", err
	}

	// Get the IP from the lease table.
	leaseip := regexp.MustCompile(`^lease (.+?) {$`)
	// Get the lease end date time.
	leaseend := regexp.MustCompile(`^\s*ends \d (.+?);$`)
	// Get the MAC address associated.
	leasemac := regexp.MustCompile(`^\s*hardware ethernet (.+?);$`)

	for _, line := range strings.Split(string(dhcpcontent), "\n") {

		if matches := leaseip.FindStringSubmatch(line); matches != nil {
			lastipmatch = matches[1]
			continue
		}

		if matches := leaseend.FindStringSubmatch(line); matches != nil {
			lastleaseendtime, _ = time.Parse("2006/01/02 15:04:05", matches[1])
			continue
		}

		if matches := leasemac.FindStringSubmatch(line); matches != nil && matches[1] == macaddr && currentleadeendtime.Before(lastleaseendtime) {
			currentip = lastipmatch
			currentleadeendtime = lastleaseendtime
		}
	}

	if currentip == "" {
		return "", fmt.Errorf("IP not found for MAC %s in DHCP leases", macaddr)
	}

	log.Debugf("IP found in DHCP lease table: %s", currentip)

	return currentip, nil
}
