/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"

	"github.com/Nanocloud/nano"
	"gopkg.in/ldap.v2"
)

var module nano.Module

var conf struct {
	Username   string
	Password   string
	ServerURL  string
	Ou         string
	LDAPServer url.URL
}

type hash map[string]interface{}

type AccountParams struct {
	UserEmail string
	Password  string
}

type ChangePasswordParams struct {
	SamAccountName string
	NewPassword    string
}

// Strucutre used in messages from RabbitMQ
type Message struct {
	Method    string
	Name      string
	Email     string
	Activated string
	Sam       string
	Password  string
}

// Plugin structure
type Ldap struct{}

type ldap_conf struct {
	host   string
	login  string
	passwd string
	ou     string
}

// Strucutre used in return messages sent to RabbitMQ
type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func DialandBind() (*ldap.Conn, error) {
	ldapConnection, err := ldap.DialTLS("tcp", conf.LDAPServer.Host,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, err
	}

	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return nil, err
	}
	return ldapConnection, nil
}

func listUsers(req nano.Request) (*nano.Response, error) {
	ldapConnection, err := DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return nil, errors.New("Search error: " + err.Error())
	}
	// Struct needed for JSON encoding
	var res struct {
		Count int
		Users []map[string]string
	}
	res.Count = len(sr.Entries)
	res.Users = make([]map[string]string, res.Count)
	i := 0

	for _, entry := range sr.Entries {
		res.Users[i] = make(map[string]string, 6)
		res.Users[i]["dn"] = entry.DN
		res.Users[i]["cn"] = entry.GetAttributeValue("cn")
		res.Users[i]["mail"] = entry.GetAttributeValue("mail")
		res.Users[i]["samaccountname"] = entry.GetAttributeValue("sAMAccountName")
		res.Users[i]["useraccountcontrol"] = entry.GetAttributeValue("userAccountControl")
		h, _ := strconv.Atoi(res.Users[i]["useraccountcontrol"])
		if h&0x0002 == 0 { // 0x0002 activated means user is disabled
			res.Users[i]["status"] = "Enabled"
		} else {
			res.Users[i]["status"] = "Disabled"

		}

		i++
	}

	return nano.JSONResponse(200, res), nil
}

func test_password(pass string) bool {
	// Windows AD password needs at leat 7 characters password,  and must contain characters from three of the following five categories:
	// uppercase character
	// lowercase character
	// digit character
	// nonalphanumeric characters
	// any Unicode character that is categorized as an alphabetic character but is not uppercase or lowercase
	if len(pass) < 7 {
		return false
	}
	d := 0
	l := 0
	u := 0
	p := 0
	o := 0
	for _, c := range pass {
		if unicode.IsDigit(c) { // check digit character
			d = 1
		} else if unicode.IsLower(c) { // check lowercase character
			l = 1
		} else if unicode.IsUpper(c) { // check uppercase character
			u = 1
		} else if unicode.IsPunct(c) { // check nonalphanumeric character
			p = 1
		} else { // other unicode character
			o = 1
		}
	}
	if d+l+u+p+o < 3 {
		return false
	}
	return true
}

// Checks if there is at least one sam account available, to use it to create a new user instead of generating a new sam account
func checkSamAvailability(ldapConnection *ldap.Conn) (error, string, int) {
	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		module.Log.Error("Search error:  " + err.Error())
		return errors.New("Search error: " + err.Error()), "", 0
	}
	count := len(sr.Entries)
	cn := ""
	for _, entry := range sr.Entries {
		h, err := strconv.Atoi(entry.GetAttributeValue("userAccountControl"))
		if err != nil {
			return errors.New("Atoi conversion error: " + err.Error()), "", 0
		}
		if h&0x0002 == 0 { //0x0002 means disabled account
		} else {
			cn = entry.GetAttributeValue("cn")
			break
		}
	}
	return nil, cn, count
}

func createNewUser(conf2 ldap_conf, params AccountParams, count int, ldapConnection *ldap.Conn) (*nano.Response, error) {
	var sam string
	if !test_password(params.Password) {
		return nano.JSONResponse(400, hash{
			"error": "Password does not meet minimum requirements",
		}), nil
	}

	dn := "cn=" + fmt.Sprintf("%d", count+1) + "," + conf2.ou

	req := ldap.NewAddRequest(dn)
	req.Attribute("objectclass", []string{"top", "person", "organizationalPerson", "User"})
	req.Attribute("mail", []string{params.UserEmail})
	pwd := encodePassword(params.Password)
	req.Attribute("unicodePwd", []string{string(pwd)})
	req.Attribute("userAccountControl", []string{"512"})
	err := ldapConnection.Add(req)
	if err != nil {
		module.Log.Error("Adding error:  " + err.Error())
		return nil, errors.New("Adding a user failed: " + err.Error())
	}

	ldapConnection, err = DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+fmt.Sprintf("%d", count+1)+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		module.Log.Error("Search error:  " + err.Error())
		return nil, errors.New("Search error: " + err.Error())
	}
	for _, entry := range sr.Entries {
		module.Log.Info(entry.GetAttributeValue("sAMAccountName"))
		sam = entry.GetAttributeValue("sAMAccountName")
	}
	return nano.JSONResponse(201, hash{
		"sam": sam,
	}), nil
}

func encodePassword(pass string) []byte {
	s := pass
	// Windows AD needs a UTF16-LE encoded password, with double quotes at the beginning and at the end
	enc := utf16.Encode([]rune(s))
	pwd := make([]byte, len(enc)*2+4)

	pwd[0] = '"'
	i := 2
	for _, n := range enc {
		pwd[i] = byte(n)
		pwd[i+1] = byte(n >> 8)
		i += 2
	}
	pwd[i] = '"'
	return pwd
}

// Uses a deactivated sam account to create a new user with it
func recycleSam(params AccountParams, ldapConnection *ldap.Conn, cn string) (*nano.Response, error) {
	var sam string
	pwd := encodePassword(params.Password)
	modify := ldap.NewModifyRequest("cn=" + cn + "," + conf.Ou)
	modify.Replace("unicodePwd", []string{string(pwd)})
	modify.Replace("userAccountControl", []string{"512"})
	modify.Replace("mail", []string{params.UserEmail})
	err := ldapConnection.Modify(modify)
	if err != nil {
		return nil, errors.New("Modify error: " + err.Error())
	}

	ldapConnection, err = DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+cn+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		module.Log.Error("Search error:  " + err.Error())
		return nil, errors.New("Search error: " + err.Error())
	}
	for _, entry := range sr.Entries {
		module.Log.Info(entry.GetAttributeValue("sAMAccountName"))
		sam = entry.GetAttributeValue("sAMAccountName")
	}
	return nano.JSONResponse(201, hash{
		"sam": sam,
	}), nil
}

func updatePassword(req nano.Request) (*nano.Response, error) {
	var params struct {
		UserEmail string
		Password  string
	}
	err := json.Unmarshal(req.Body, &params)
	if err != nil {
		module.Log.Error("Unable to unmarshall params: " + err.Error())
		return nil, err
	}

	if len(req.Params["user_id"]) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "user id is missing",
		}), nil
	}

	params.UserEmail = req.Params["user_id"]

	ldapConnection, err := DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	defer ldapConnection.Close()

	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+params.UserEmail+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		module.Log.Error("Searching error: " + err.Error())
		return nil, errors.New("Search error: " + err.Error())
	}

	var cn string
	if len(sr.Entries) != 1 {
		module.Log.Error("Invalid Email")
		return nil, errors.New("invalid Email")
	}
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")

	}
	pwd := encodePassword(params.Password)

	modify := ldap.NewModifyRequest("cn=" + cn + "," + conf.Ou)
	modify.Replace("unicodePwd", []string{string(pwd)})
	err = ldapConnection.Modify(modify)
	if err != nil {
		module.Log.Error("Password modification failed: " + err.Error())
		return nil, errors.New("Password modification failed: " + err.Error())
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func createUser(req nano.Request) (*nano.Response, error) {
	var params struct {
		UserEmail string
		Password  string
	}

	err := json.Unmarshal(req.Body, &params)

	if err != nil {
		module.Log.Error("Unable to unmarshall params: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	if params.UserEmail == "" || params.Password == "" {
		module.Log.Error("Email or password missing")
		return nano.JSONResponse(400, hash{
			"error": "Email or password missing",
		}), nil
	}

	var tconf ldap_conf
	tconf.host = conf.LDAPServer.Scheme + "://" + conf.LDAPServer.Host
	tconf.login = conf.Username
	tconf.passwd = conf.Password
	tconf.ou = conf.Ou
	// return "", to ldap go API to set the password

	ldapConnection, err := DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	defer ldapConnection.Close()

	err, cn, count := checkSamAvailability(ldapConnection) // if an account is disabled, this function will look for his CN
	if err != nil {
		module.Log.Error("Error while checking sam availability: " + err.Error())
		return nil, err
	}

	// if no disabled accounts were found, a real new user is created
	if cn == "" {
		res, err := createNewUser(tconf, params, count, ldapConnection)
		if err != nil {
			module.Log.Error(err.Error())
			return nil, err
		}
		return res, nil
	}

	// if a disabled account is found, modifying this account instead of creating a new one
	return recycleSam(params, ldapConnection, cn)
}

func forcedisableAccount(req nano.Request) (*nano.Response, error) {
	userId := req.Params["user_id"]

	if len(userId) < 1 {
		module.Log.Error("User ID missing")
		return nano.JSONResponse(400, hash{
			"error": "User id is missing",
		}), nil
	}

	ldapConnection, err := DialandBind()
	if err != nil {
		module.Log.Error("Error while connection to Active Directory: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		conf.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+userId+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		module.Log.Error("Searching error: " + err.Error())
		return nil, errors.New("Searching error: " + err.Error())
	}

	if len(sr.Entries) != 1 {
		module.Log.Error("Email does not match any user, or several users have the same mail adress")
		// means entered mail was not valid, or several user have the same mail
		return nano.JSONResponse(400, hash{
			"error": "Email does not match any user, or several users have the same mail adress",
		}), nil
	}
	var cn string
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")
	}
	modify := ldap.NewModifyRequest("cn=" + cn + "," + conf.Ou)
	modify.Replace("userAccountControl", []string{"514"}) // 512 is a normal account, 514 is disabled ( 512 + 0x0002 )
	err = ldapConnection.Modify(modify)
	if err != nil {
		module.Log.Error("Modify  error: " + err.Error())
		return nil, errors.New("Modify error: " + err.Error())
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func getCert(server url.URL) error {
	err := os.MkdirAll("/opt/conf", 755)
	if err != nil {
		return err
	}

	if server.User == nil {
		return errors.New("No authentication informations provided")
	}

	hostname := strings.SplitN(server.Host, ":", 2)

	password, ok := server.User.Password()
	if !ok {
		return errors.New("Password not set")
	}
	cmd := exec.Command(
		"sshpass", "-p", password,
		"scp", "-P", hostname[1],
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf(
			"%s@%s:%s",
			server.User.Username(),
			hostname[0],
			"/cygdrive/c/users/administrator/ad2012.cer",
		),
		"/opt/conf",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		module.Log.Error("Failed to execute script ", err, string(response))
		return err
	}
	return nil
}

func genLdaprc(host string) error {
	err := os.MkdirAll("/etc/ldap", 755)
	if err != nil {
		return err
	}

	base := env("BASE", "DC=intra,DC=localdomain,DC=com")
	bindDn := env("BIND_DN", "CN=Administrator,DC=intra,DC=localdomain,DC=com")
	tlsCacert := env("TLS_CACERT", "/opt/conf/ad2012.cer")

	s := fmt.Sprintf(`
BASE %s
BINDDN %s
URI %s
TLS_CACERT %s
TLS_REQCERT never
`,
		base,
		bindDn,
		host,
		tlsCacert,
	)

	return ioutil.WriteFile("/etc/ldap/ldaprc", []byte(s), 0644)
}

func main() {
	conf.Username = env("LDAP_USERNAME", "CN=Administrator,CN=Users,DC=intra,DC=localdomain,DC=com")
	conf.Password = env("LDAP_PASSWORD", "Nanocloud123+")
	conf.Ou = env("LDAP_OU", "OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")

	ldapServer, err := url.Parse(env("LDAP_SERVER_URI", "ldaps://Administrator:Nanocloud123+@172.17.0.1:6003"))
	if err != nil {
		panic(err)
	}

	sshServer, err := url.Parse(env("SSH_SERVER_URI", "ssh://Administrator:Nanocloud123+@172.17.0.1:6001"))
	if err != nil {
		panic(err)
	}

	conf.LDAPServer = *ldapServer

	module = nano.RegisterModule("ldap")

	genLdaprc(ldapServer.Host)

	try := 0
	for try = 0; try < 10; try++ {
		err := getCert(*sshServer)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}

	if try == 10 {
		module.Log.Fatal("Unable to connect to windows")
	}

	module.Post("/ldap/users", createUser)
	module.Get("/ldap/users", listUsers)
	module.Put("/ldap/users/:user_id", updatePassword)
	module.Post("/ldap/users/:user_id/disable", forcedisableAccount)

	module.Listen()
}
