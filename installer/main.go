/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
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

package main

/*
#include <sys/sysinfo.h>
*/
import (
	"C"
)

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Sysinfo_t C.struct_sysinfo

var (
	prefix        string         = "/var/lib/nanocloud"
	IaasApiUrl    string         = "http://localhost:8082/"
	cpuinfoRegExp *regexp.Regexp = regexp.MustCompile("([^:]*?)\\s*:\\s*(.*)$")
	sysinfo       C.struct_sysinfo
)

func checkPrerequisite() {

	// Check if user is root
	user, err := user.Current()
	if err != nil || user.Uid != "0" {
		log.Fatalf("This script must be run as root")
		os.Exit(1)
	}

	// Check system total memory
	ret := C.sysinfo(&sysinfo)
	if int(ret) == -1 || int(sysinfo.totalram) <= 2*1024*1024*1024 {
		log.Fatalf("Not enough RAM available")
		os.Exit(2)
	}

	// Check intel VT is enabled
	// Information can be found here :
	// https://www.redhat.com/archives/libvir-list/2007-March/msg00218.html
	b, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Fatalf("Cannot read /proc/cpuinfo")
	}

	content := string(b)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		var key string
		var value string

		if len(line) == 0 && i != len(lines)-1 {
			continue
		} else if i == len(lines)-1 {
			continue
		}

		submatches := cpuinfoRegExp.FindStringSubmatch(line)
		key = submatches[1]
		value = submatches[2]

		if key != "flags" {
			continue
		}

		if strings.Count(value, "vmx") == 0 && strings.Count(value, "svm") == 0 {
			log.Fatalf("Hardware virtualization is not enable in your system")
			os.Exit(3)
		}
	}
}

func unpackFile() {
	var (
		splitAssetName []string
		cannonicalName string
		path           string
		mode           os.FileMode = 0755
	)

	os.MkdirAll(prefix, mode)

	for _, AssetName := range AssetNames() {
		splitAssetName = strings.Split(AssetName, "/")
		path = strings.Join(splitAssetName[1:len(splitAssetName)-1], "/")
		cannonicalName = splitAssetName[len(splitAssetName)-1]

		path = filepath.Join(prefix, path)

		os.MkdirAll(path, mode)

		unpackAsset(AssetName, filepath.Join(path, cannonicalName))
	}

	extractTarGzFile(filepath.Join(prefix, "iaas.tar.gz"))
}

func startIaasAPI() {

	cmd := exec.Command("/etc/init.d/iaasAPI", "start")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to start iaasAPI, error: %s, output: %s", err, string(out))
	}
}

func lauchVM(VMName string) {
	var params = map[string]string{"name": VMName}
	log.Println("DEBUG : params ", params)
	jsonRpcRequest(
		IaasApiUrl,
		"Iaas.Start",
		params,
	)
}

func domainLookup() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	// TODO How to ensure the addrs returned is good ?
	return addrs[len(addrs)-1]
}

func downloadFromUrl(url string, retry int) bool {
	var success bool = false

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for attempt := 1; attempt <= retry; attempt++ {
		response, err := client.Get(url)
		if err == nil {
			success = true
			defer response.Body.Close()
		}
		if success {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	return success
}

func launchAPI() {
	log.Println("Starting API…")
	startIaasAPI()
	log.Println("  NOTES: To automatically launch the API, run the following command as root:")
	log.Println("  # update-rc.d \"iaasAPI\" defaults")
	log.Println("  API OK")
}

func launch() {
	log.Println("Booting first VM…")
	//lauchVM("noauth-mini-free_use-10.104.16.190-linux-alpine-3.2-x86_64")
	lauchVM("coreos")
	log.Println("  boot first VM OK")

	// TODO Hardcoded IP and port number
	if downloadFromUrl("https://localhost:2224/login.html", 30) {
		log.Println("Setup complete")
		log.Printf("You can now manage your platform on : https://localhost:2224\n")
		log.Println("Default admin credential:")
		log.Println("\tEmail: admin@nanocloud.com")
		log.Println("\tPassword: admin")
		log.Println("This URL will only be accessible from this host.")
		log.Println("If you want to access it from another machine : https://<public-ip-or-fqdn>:8490\n")
		log.Println("")
		log.Println("Use the following commands as root to start, stop or get status information")
		log.Printf("    # %s\n", filepath.Join(prefix, "scripts/start.sh"))
		log.Printf("    # %s\n", filepath.Join(prefix, "scripts/stop.sh"))
		log.Printf("    # %s\n", filepath.Join(prefix, "scripts/status.sh"))
	}
}

func createSymlink() {
	apiInitScriptPath := filepath.Join(prefix, "scripts/APIinitScript")

	bExists := false

	if _, err := os.Lstat("/etc/init.d/iaasAPI"); err != nil {
		bExists = os.IsExist(err)
	} else {
		log.Println("Exists")
		bExists = true
	}

	if bExists {
		if err := os.Remove("/etc/init.d/iaasAPI"); err != nil {
			log.Fatalf("Unable to remove symlink existing before installation : /etc/init.d/iaasAPI with error: %s", err)
		}
	}

	err := os.Symlink(apiInitScriptPath, "/etc/init.d/iaasAPI")
	if err != nil {
		log.Fatalf("Error creating symlink /etc/init.d/iaasAPI : %s", err)
	}
}

func installFromBundledData() {
	log.Println("Unpacking files…")
	unpackFile()
	log.Println("  unpacking OK")
}

func installForExistingBuild() {
	log.Println("Install for existing build…")
	var (
		mode os.FileMode = 0755
	)
	os.MkdirAll(prefix, mode)
}

func launchCheckPrerequisite() {
	log.Println("Checking prerequisite…")
	checkPrerequisite()
	log.Println("  prerequisite OK")
}

func main() {
	var installWillProceedMsg string = `Installation will proceed…
    To uninstall Nanocloud community, run the followind command:
    curl https://community.nanocloud.com/nanocloud_uninstall.sh | sh`

	if len(os.Args) == 1 {
		launchCheckPrerequisite()
		log.Println(installWillProceedMsg)
		installFromBundledData()
		createSymlink()
		launchAPI()
	} else if len(os.Args) > 1 {
		switch os.Args[1] {
		case "release":
			launchCheckPrerequisite()
			log.Println(installWillProceedMsg)
			installForExistingBuild()
			createSymlink()
			launchAPI()
		case "unpack":
			launchCheckPrerequisite()
			log.Println(installWillProceedMsg)
			unpackFile()
			createSymlink()
			launchAPI()
		case "launch":
			launch()
		}
	}
}
