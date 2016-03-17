package provisioning

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Nanocloud/community/plaza/server/router"
	log "github.com/Sirupsen/logrus"

	//"github.com/labstack/echo"
)

type hash map[string]interface{}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("powershell.exe", command)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return string(resp), err
	}
	return string(resp), nil
}

var commands = [][]string{
	{
		"New-Item HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows -Name WindowsUpdate",
		"New-Item HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate -Name AU",
		"New-ItemProperty HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate\\AU -Name NoAutoUpdate -Value 1",
	},
	{
		"Install-windowsfeature AD-domain-services",
		"Import-Module ADDSDeployment; $pwd=ConvertTo-SecureString 'Nanocloud123+' -asplaintext -force; Install-ADDSForest -CreateDnsDelegation:$false -DatabasePath 'C:\\Windows\\NTDS' -DomainMode 'Win2012R2' -DomainName 'intra.localdomain.com' -SafeModeAdministratorPassword:$pwd -DomainNetbiosName 'INTRA' -ForestMode 'Win2012R2' -InstallDns:$true -LogPath 'C:\\Windows\\NTDS' -NoRebootOnCompletion:$true -SysvolPath 'C:\\Windows\\SYSVOL' -Force:$true",
	},
	{
		"set-ItemProperty -Path 'HKLM:\\System\\CurrentControlSet\\Control\\Terminal Server'-name 'fDenyTSConnections' -Value 0",
		"Enable-NetFirewallRule -DisplayGroup 'Remote Desktop'",
		"set-ItemProperty -Path 'HKLM:\\System\\CurrentControlSet\\Control\\Terminal Server\\WinStations\\RDP-Tcp' -name 'UserAuthentication' -Value 1",
	},
	{
		"import-module RemoteDesktop; Import-module ServerManager; Add-WindowsFeature -Name RDS-RD-Server -IncludeAllSubFeature; Add-WindowsFeature -Name RDS-Web-Access -IncludeAllSubFeature; Add-WindowsFeature -Name RDS-Connection-Broker -IncludeAllSubFeature; Install-windowsfeature RSAT-AD-AdminCenter",
	},
	{
		"NEW-ADOrganizationalUnit 'NanocloudUsers' -path 'DC=intra,DC=localdomain,DC=com'",
	},
	{
		"Import-Module ServerManager; Add-WindowsFeature Adcs-Cert-Authority",
		"$secpasswd = ConvertTo-SecureString 'Nanocloud123+' -AsPlainText -Force;$mycreds = New-Object System.Management.Automation.PSCredential ('Administrator', $secpasswd); Install-AdcsCertificationAuthority -CAType 'EnterpriseRootCa' -Credential:$mycreds -force:$true ",
	},
	{
		"sc.exe config RDMS start= auto",
		"New-NetFirewallRule -Protocol TCP -LocalPort 9090 -Direction Inbound -Action Allow -DisplayName PLAZA",
		"import-module remotedesktop ; New-RDSessionDeployment -ConnectionBroker adapps.intra.localdomain.com -WebAccessServer adapps.intra.localdomain.com -SessionHost adapps.intra.localdomain.com; New-RDSessionCollection -CollectionName collection -SessionHost adapps.intra.localdomain.com -CollectionDescription 'Nanocloud collection' -ConnectionBroker adapps.intra.localdomain.com",
		//		"import-module remotedesktop ;New-RDRemoteApp -CollectionName collection -DisplayName hapticPowershell -FilePath 'C:\\Windows\\system32\\WindowsPowerShell\\v1.0\\powershell.exe' -Alias hapticPowershell -CommandLineSetting Require -RequiredCommandLine '-ExecutionPolicy Bypass c:\\publishApplication.ps1'",
	},
}

var confPath = "C:/prov.json"

func executeCommands(commands []string) error {
	var err error
	var resp string
	for _, cmd := range commands {
		resp, err = executeCommand(cmd)
		log.Error(string(resp))
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func markAsDone(conf []bool, action int) {
	conf[action] = true
	b, err := json.Marshal(conf)
	if err != nil {
		log.Error(err)
		return
	}
	err = ioutil.WriteFile(confPath, b, 0644)
}

func LaunchAll() {
	go router.Start()
	ProvisionAll()
}

func waitForADWS() {
	for {
		resp, err := executeCommand("Write-Host (Get-Service -Name ADWS).status")
		log.Error(string(resp))
		if err != nil {
			log.Error(err)
		}
		if strings.Contains(string(resp), "Running") {
			break
		}
		log.Error("Waiting for ADWS to be running")
		time.Sleep(time.Second * 5)
	}
}

func ProvisionAll() {
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		os.Create(confPath)
		var p = make([]bool, 7)
		b, err := json.Marshal(p)
		if err != nil {
			log.Error(err)
		}
		err = ioutil.WriteFile(confPath, b, 0644)
	}

	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Error(err)
		return
	}
	var conf []bool
	err = json.Unmarshal(file, &conf)
	for index, done := range conf {
		if !done {
			for i, val := range commands[index:] {
				log.Error("NEW PACK")
				if i+index == 4 {
					waitForADWS()
				}
				err = executeCommands(val)
				if err == nil {
					markAsDone(conf, i+index)
				} else {
					log.Error("THERE WAS AN ERROR")
				}
				if i+index == 3 {
					executeCommand("Restart-Computer -Force")
					os.Exit(0)
				}
			}
			break
		}
	}
}

/*
func DisableWU(c *echo.Context) error {
	return executeCommands(commands["disablewu"], c)
}

func reterr(e error, resp string, c *echo.Context) error {
	return c.JSON(
		http.StatusInternalServerError,
		hash{
			"error": []hash{
				hash{
					"title":  e.Error(),
					"detail": resp,
				},
			},
		},
	)
}

func retok(c *echo.Context) error {
	return c.JSON(
		http.StatusOK,
		hash{
			"data": hash{
				"success": true,
			},
		},
	)
}

func CheckWU(c *echo.Context) error {
	resp, err := executeCommand("Get-ItemProperty HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\WindowsUpdate\\AU")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "NoAutoUpdate : 1") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"windows-update": "disabled",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"windows-update": "enabled",
				},
			},
		)
	}
}

func InstallAD(c *echo.Context) error {
	return executeCommands(commands["installad"], c)
}

func CheckAD(c *echo.Context) error {
	resp, err := executeCommand("Get-ADForest")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "intra.localdomain.com") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "Nanocloud forest installed",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "Nanocloud forest not installed",
				},
			},
		)
	}
}

func EnableRDP(c *echo.Context) error {
	return executeCommands(commands["enablerdp"], c)
}

func CheckRDP(c *echo.Context) error {
	resp, err := executeCommand("Write-Host (Get-Service -Name RDMS).status")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "Running") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "RDP Service running",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "RDP Service is down",
				},
			},
		)
	}
}

func InstallRDS(c *echo.Context) error {
	return executeCommands(commands["installrds"], c)
}

func CheckRDS(c *echo.Context) error {
	resp, err := executeCommand("Write-Host (Get-Service -Name TermService).status")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "Running") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "RDS Service running",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "RDS Service is down",
				},
			},
		)
	}
}

func CreateOU(c *echo.Context) error {
	return executeCommands(commands["createou"], c)
}

func CheckOU(c *echo.Context) error {
	resp, err := executeCommand("Get-ADOrganizationalUnit -Filter 'Name -like \"NanocloudUsers\"'")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "NanocloudUsers") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"organizational-unit": "created",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"organizational-unit": "Not created",
				},
			},
		)
	}
}

func InstallADCS(c *echo.Context) error {
	return executeCommands(commands["installadcs"], c)
}

func CheckADCS(c *echo.Context) error {
	resp, err := executeCommand("Write-Host (Get-Service -Name CertSvc).status")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "Running") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "ADCS Service running",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"state": "ADCS Service is down",
				},
			},
		)
	}
}

func SessionDeploy(c *echo.Context) error {
	return executeCommands(commands["sessiondeploy"], c)
}

func CheckCollection(c *echo.Context) error {
	resp, err := executeCommand("import-module remotedesktop; Get-RDSessionCollection -CollectionName 'collection'")
	if err != nil {
		return reterr(err, resp, c)
	}
	if strings.Contains(resp, "collection") {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"status": "collection created",
				},
			},
		)
	} else {
		return c.JSON(
			http.StatusOK,
			hash{
				"data": hash{
					"status": "collection not created",
				},
			},
		)
	}
}*/
