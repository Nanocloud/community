/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2016 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package plaza

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/provisioner"
	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type Cmd_t struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Domain   string   `json:"domain"`
	Command  []string `json:"command"`
	Stdin    string   `json:"stdin"`
}

type result_t struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   string `json:"code"`
}

func Exec(address string, port int, cmd *Cmd_t) (*result_t, error) {
	var err error

	client := &http.Client{}

	instr, err := json.Marshal(&cmd)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(instr)

	var resp *http.Response

	for i := 0; i < 10; i++ {
		resp, err = client.Post(
			fmt.Sprintf("http://%s:%d/exec", address, port),
			"application/json",
			buff,
		)
		if err == nil {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	r := result_t{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func PowershellExec(
	address string, port int,
	username string, domain string, password string,
	command ...string,
) (*result_t, error) {
	cmd := Cmd_t{
		Username: username,
		Password: password,
		Domain:   domain,
		Command: []string{
			"C:\\Windows\\System32\\WindowsPowershell\\v1.0\\powershell.exe",
			"-Command",
			"-",
		},
		Stdin: strings.Join(command, " "),
	}

	res, err := Exec(address, port, &cmd)
	if err != nil {
		return nil, err
	}

	if res.Code != "exit status 0" {
		return nil, errors.New("STDOUT: " + res.Stdout + "\nSTDERR: " + res.Stderr)
	}
	return res, nil
}

func PublishApp(
	address string, port int,
	username string, domain string, password string,
	collectionName string, displayName string, filePath string,
) ([]byte, error) {
	res, err := PowershellExec(
		address, port,
		username, domain, password,
		"Try {",
		"Import-module RemoteDesktop;",
		fmt.Sprintf(
			"New-RDRemoteApp -CollectionName '%s' -DisplayName '%s' -FilePath '%s' -ErrorAction Stop | ConvertTo-Json",
			collectionName, displayName, filePath,
		),
		"}",
		"Catch {",
		"$ErrorMessage = $_.Exception.Message;",
		"Write-Output -InputObject $ErrorMessage;",
		"exit 1;",
		"}",
	)

	if err != nil {
		return nil, err
	}

	out := res.Stdout

	return []byte(out), nil
}

func UnpublishApp(
	address string, port int,
	username string, domain string, password string,
	collectionName string, alias string,
) ([]byte, error) {
	res, err := PowershellExec(
		address, port,
		username, domain, password,
		"Try {",
		"Import-module RemoteDesktop;",
		fmt.Sprintf(
			"Remove-RDRemoteApp -CollectionName '%s' -Alias '%s' -Force -ErrorAction Stop | ConvertTo-Json",
			collectionName, alias,
		),
		"}",
		"Catch {",
		"$ErrorMessage = $_.Exception.Message;",
		"Write-Output -InputObject $ErrorMessage;",
		"exit 1;",
		"}",
	)

	if err != nil {
		return nil, err
	}

	out := res.Stdout

	return []byte(out), nil
}

func checkPlaza(ip string, port string) bool {
	_, err := http.Get("http://" + ip + ":" + port + "/")
	if err != nil {
		return false
	} else {
		return true
	}
}

func provExec(p io.Writer, machine vms.Machine, command string) (string, error) {

	machine.Status()
	ip, _ := machine.IP()
	for checkPlaza(ip.String(), utils.Env("PLAZA_PORT", "9090")) != true {

		machine.Status()
		ip, _ = machine.IP()
		time.Sleep(time.Millisecond * 500)
	}

	plazaAddress, err := machine.IP()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	plazaPort, err := strconv.Atoi(utils.Env("PLAZA_PORT", "9090"))
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	domain := utils.Env("WINDOWS_DOMAIN", "")
	if domain == "" {
		log.Error("domain unknown")
		return "", errors.New("domain unknown")
	}

	username, password, err := machine.Credentials()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	res, err := PowershellExec(
		plazaAddress.String(),
		plazaPort,
		username,
		domain,
		password,
		command,
	)

	p.Write([]byte(command))
	if err != nil {
		return "", err
	}

	return res.Stdout, nil
}

func isStopped(machine vms.Machine) bool {

	status, err := machine.Status()
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if status == vms.StatusDown {
		return true
	} else {
		return false
	}
}

func Provision(machine vms.Machine) provisioner.ProvFunc {

	return func(p io.Writer) {
		resp, err := provExec(p, machine, "New-Item HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows -Name WindowsUpdate")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		username, password, err := machine.Credentials()
		if err != nil {
			p.Write([]byte(err.Error()))
		}
		pcname, err := provExec(p, machine, "hostname")
		if err != nil {
			p.Write([]byte(err.Error()))
		}
		pcname = strings.TrimSpace(pcname)
		domain := utils.Env("WINDOWS_DOMAIN", "")

		resp, err = provExec(p, machine, "New-Item HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate -Name AU")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "New-ItemProperty HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate\\AU -Name NoAutoUpdate -Value 1")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "Install-windowsfeature AD-domain-services")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "Import-Module ADDSDeployment; $pwd=ConvertTo-SecureString '"+password+"' -asplaintext -force; Install-ADDSForest -CreateDnsDelegation:$false -DatabasePath 'C:\\Windows\\NTDS' -DomainMode 'Win2012R2' -DomainName '"+domain+"' -SafeModeAdministratorPassword:$pwd -DomainNetbiosName 'INTRA' -ForestMode 'Win2012R2' -InstallDns:$true -LogPath 'C:\\Windows\\NTDS' -NoRebootOnCompletion:$true -SysvolPath 'C:\\Windows\\SYSVOL' -Force:$true")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		machine.Stop()
		for isStopped(machine) != true {
		}
		machine.Start()

		resp, err = provExec(p, machine, "set-ItemProperty -Path 'HKLM:\\System\\CurrentControlSet\\Control\\Terminal Server'-name 'fDenyTSConnections' -Value 0")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "Enable-NetFirewallRule -DisplayGroup 'Remote Desktop'")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "set-ItemProperty -Path 'HKLM:\\System\\CurrentControlSet\\Control\\Terminal Server\\WinStations\\RDP-Tcp' -name 'UserAuthentication' -Value 1")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "import-module RemoteDesktop; Import-module ServerManager; Add-WindowsFeature -Name RDS-RD-Server -IncludeAllSubFeature; Add-WindowsFeature -Name RDS-Web-Access -IncludeAllSubFeature; Add-WindowsFeature -Name RDS-Connection-Broker -IncludeAllSubFeature")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "import-module RemoteDesktop; Import-module ServerManager; Install-windowsfeature RSAT-AD-AdminCenter")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		machine.Stop()
		for isStopped(machine) != true {

		}
		machine.Start()

		resp, err = provExec(p, machine, "sc.exe config RDMS start= auto")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "Import-Module ServerManager; Add-WindowsFeature Adcs-Cert-Authority")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "$secpasswd = ConvertTo-SecureString '"+password+"' -AsPlainText -Force;$mycreds = New-Object System.Management.Automation.PSCredential ('"+username+"', $secpasswd); Install-AdcsCertificationAuthority -CAType 'EnterpriseRootCa' -Credential:$mycreds -force:$true ")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "Start-Service RDMS; import-module remotedesktop ; New-RDSessionDeployment -ConnectionBroker "+pcname+"."+domain+" -WebAccessServer "+pcname+"."+domain+" -SessionHost "+pcname+"."+domain)
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		time.Sleep(time.Second * 60)

		resp, err = provExec(p, machine, "import-module remotedesktop ; New-RDSessionCollection -CollectionName collection -SessionHost "+pcname+"."+domain+" -CollectionDescription 'Nanocloud collection' -ConnectionBroker "+pcname+"."+domain)
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "(Get-WmiObject -class 'Win32_TSGeneralSetting' -Namespace root\\cimv2\\terminalservices -ComputerName "+pcname+").SetUserAuthenticationRequired(0)")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}

		resp, err = provExec(p, machine, "NEW-ADOrganizationalUnit 'NanocloudUsers' -path 'DC=intra,DC=localdomain,DC=com'")
		if err != nil {
			p.Write([]byte(err.Error()))
		} else {
			p.Write([]byte(resp))
		}
	}
}
