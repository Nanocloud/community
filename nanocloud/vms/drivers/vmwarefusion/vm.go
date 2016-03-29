package vmwarefusion

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type vm struct {
	plazaLocation string
	storageDir    string
}

func (v *vm) createMachineStorage(vmId string, t *machineType) (string, error) {
	log.WithFields(log.Fields{
		"VM": vmId,
	}).Info("Create VM storage")

	dstDir := path.Join(v.storageDir, path.Join("vm", vmId))
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return "", err
	}

	dst := path.Join(dstDir, "disk.vmdk")

	log.Debugln(
		"Executing:",
		vdiskmanager, "-c", "-t", "0",
		"-s", t.size,
		"-a", "lsilogic",
		dst,
	)
	cmd := exec.Command(
		vdiskmanager, "-c", "-t", "0",
		"-s", t.size,
		"-a", "lsilogic",
		dst,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{
			"VM": vmId,
		}).Error(string(out))
		return "", err
	}

	return dst, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	vmx := path.Join(v.storageDir, "vm", id, "conf.vmx")

	_, err := os.Stat(vmx)
	if err != nil {
		return nil, err
	}

	return &machine{
		id:         id,
		storageDir: v.storageDir,
	}, nil
}

func (v *vm) Machines() ([]vms.Machine, error) {
	vmDir, err := os.Open(path.Join(v.storageDir, "vm"))
	if err != nil {
		return nil, err
	}

	ids, err := vmDir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	machines := make([]vms.Machine, 0)
	for _, id := range ids {
		if string(id[0]) == "." {
			continue
		}
		machines = append(machines, &machine{
			id:         id,
			storageDir: v.storageDir,
		})
	}

	return machines, nil
}

func (v *vm) createVmxForSetup(vmId string, name string, t *machineType, hdd string, iso string, installISO string) (string, error) {
	log.WithFields(log.Fields{
		"VM": vmId,
	}).Info("Create setup VMX file")

	dstDir := path.Join(v.storageDir, "vm", vmId)
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return "", err
	}

	vmx := template.Must(template.New("vmx").Parse(vmxTemplate))

	dst := path.Join(dstDir, "conf.vmx")
	vmxFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}

	var conf struct {
		Name                string
		WindowsInstallISO   string
		NanocloudInstallISO string
		RAM                 int
		CPU                 int
		WMDKHardDrive       string
	}

	conf.Name = name
	conf.WindowsInstallISO = iso
	conf.NanocloudInstallISO = installISO
	conf.RAM = t.ram
	conf.CPU = t.cpu
	conf.WMDKHardDrive = hdd

	err = vmx.Execute(vmxFile, &conf)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"VM": vmId,
	}).Info("Setup VMX file created")

	return dst, nil
}

func (v *vm) Create(name string, password string, rawType vms.MachineType) (vms.Machine, error) {
	h := sha1.New()
	h.Write([]byte(uuid.NewV4().String()))
	id := hex.EncodeToString(h.Sum(nil))[0:14]

	_, err := os.Stat(path.Join(v.storageDir, "vm", id))
	if err == nil {
		return nil, errors.New("A machine with the same id exists already")
	}

	if !os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"VM": id,
		}).Error(err)
		return nil, errors.New("Unable to check machine id uniqueness")
	}

	if rawType == nil {
		rawType = defaultType
	}

	t, ok := rawType.(*machineType)
	if !ok {
		return nil, errors.New("VM Type not supported")
	}

	iso, err := v.downloadWindowsISO()
	if err != nil {
		return nil, err
	}

	hdd, err := v.createMachineStorage(id, t)
	if err != nil {
		return nil, err
	}

	installISO, err := v.createInstallDisk(id, password)
	if err != nil {
		return nil, err
	}

	_, err = v.createVmxForSetup(id, name, t, hdd, iso, installISO)
	if err != nil {
		return nil, err
	}

	m := machine{
		id:         id,
		storageDir: v.storageDir,
	}

	return &m, nil
}

func (v *vm) downloadWindowsISO() (string, error) {
	dstDir := path.Join(v.storageDir, "downloads")
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return "", err
	}

	f := path.Join(dstDir, "windows.iso")
	s, err := os.Stat(f)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		log.Info("Downloading windows ISO")
		img, err := os.Create(f)
		if err != nil {
			return "", err
		}

		r, err := http.Get(windowsIsoUri)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(img, r.Body)
		log.Info("Windows ISO downloaded")
		if err != nil {
			return "", err
		}
	} else {
		log.Info("Windows ISO already downloaded")
	}

	if !s.Mode().IsRegular() {
		return f, errors.New("Windows ISO is not a file")
	}
	return f, nil
}

func (v *vm) createInstallDisk(vmId string, password string) (string, error) {
	log.WithFields(log.Fields{
		"VM": vmId,
	}).Info("Create installation disk")

	dstDir := path.Join(v.storageDir, path.Join("vm", vmId))
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return "", err
	}

	dst := path.Join(dstDir, "install.iso")

	installDir := path.Join(dstDir, "install")
	err = os.MkdirAll(installDir, 0755)
	if err != nil {
		return "", err
	}

	err = utils.CopyFile(v.plazaLocation, path.Join(installDir, "plaza.exe"))
	if err != nil {
		return "", err
	}

	autounattendTmpl := template.Must(template.New("autounattend").Parse(autounattend))

	autounattendFile, err := os.Create(path.Join(installDir, "Autounattend.xml"))
	if err != nil {
		return "", err
	}

	var conf struct {
		Hostname string
		TimeZone string
		Password string
	}

	conf.Hostname = "adapps"
	conf.TimeZone = "Central Europe Standard Time"
	conf.Password = password

	err = autounattendTmpl.Execute(autounattendFile, &conf)
	if err != nil {
		return "", err
	}

	log.Debugln(
		"Executing:", "hdiutil", "makehybrid", "-iso", "-joliet",
		"-o", dst, installDir,
	)
	cmd := exec.Command(
		"hdiutil", "makehybrid", "-iso", "-joliet",
		"-o", dst, installDir,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{
			"VM": vmId,
		}).Error(string(out))
		return "", err
	}
	log.WithFields(log.Fields{
		"VM": vmId,
	}).Info("Installation disk created")
	return dst, nil
}

func (v *vm) Types() ([]vms.MachineType, error) {
	return []vms.MachineType{defaultType}, nil
}
