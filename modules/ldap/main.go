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

/*#include <ldap.h>
#include <stdlib.h>
#include <sys/time.h>
#include <stdio.h>
#include <lber.h>
typedef struct ldapmod_str {
	int	 mod_op;
	char	  *mod_type;
	char    **mod_vals;
} LDAPModStr;
int _ldap_add(LDAP *ld, char* dn, LDAPModStr **attrs){
	return ldap_add_ext_s(ld, dn, (LDAPMod **)attrs, NULL, NULL);
}
*/
// #cgo CFLAGS: -DLDAP_DEPRECATED=1
// #cgo linux CFLAGS: -DLINUX=1
// #cgo LDFLAGS: -lldap -llber
import "C"

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
	"unsafe"

	"github.com/Nanocloud/nano"
	"gopkg.in/ldap.v2"
)

const (
	LDAP_OPT_SUCCESS          = 0
	LDAP_OPT_ERROR            = -1
	LDAP_VERSION3             = 3
	LDAP_OPT_PROTOCOL_VERSION = 0x0011
	LDAP_SUCCESS              = 0x00
	LDAP_NO_LIMIT             = 0
	LDAP_OPT_REFERRALS        = 0x0008
	LDAP_MOD_REPLACE          = 0x0002
)

const (
	LDAP_SCOPE_BASE        = 0x0000
	LDAP_SCOPE_ONELEVEL    = 0x0001
	LDAP_SCOPE_SUBTREE     = 0x0002
	LDAP_SCOPE_SUBORDINATE = 0x0003 // openLDAP extension
	LDAP_SCOPE_DEFAULT     = -1     // openLDAP extension
)

var module nano.Module

var conf struct {
	Username   string
	Password   string
	ServerURL  string
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
	ldapConnection *C.LDAP
	host           string
	login          string
	passwd         string
	ou             string
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

// Setting LDAP version and referrals
func setOptions(ldapConnection *C.LDAP) error {
	var version C.int
	var opt C.int
	version = LDAP_VERSION3
	opt = 0
	err := C.ldap_set_option(ldapConnection, LDAP_OPT_PROTOCOL_VERSION, unsafe.Pointer(&version))
	if err != LDAP_SUCCESS {
		return errors.New("Options settings error: " + C.GoString(C.ldap_err2string(err)))
	}
	err = C.ldap_set_option(ldapConnection, LDAP_OPT_REFERRALS, unsafe.Pointer(&opt))
	if err != LDAP_SUCCESS {
		return errors.New("Options settings error: " + C.GoString(C.ldap_err2string(err)))
	}
	return nil
}

func listUsers(req nano.Request) (*nano.Response, error) {
	ldapConnection, err := ldap.DialTLS("tcp", conf.LDAPServer.Host,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, errors.New("Dial error: " + err.Error())
	}

	module.Log.Info(conf.Username)
	module.Log.Info(conf.Password)

	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
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
func deleteUsers(mails []string) error {
	var tconf ldap_conf

	tconf.host = conf.ServerURL
	tconf.login = conf.Username
	tconf.passwd = conf.Password
	var version C.int
	var v C.int
	version = LDAP_VERSION3
	v = 0
	err := C.ldap_set_option(tconf.ldapConnection, LDAP_OPT_PROTOCOL_VERSION, unsafe.Pointer(&version))
	if err != LDAP_SUCCESS {
		return errors.New("Options settings error: " + C.GoString(C.ldap_err2string(err)))
	}

	err = C.ldap_set_option(tconf.ldapConnection, LDAP_OPT_REFERRALS, unsafe.Pointer(&v))
	if err != LDAP_SUCCESS {
		return errors.New("Deletion error: " + C.GoString(C.ldap_err2string(err)))
	}

	rc := C.ldap_initialize(&tconf.ldapConnection, C.CString(tconf.host))
	if tconf.ldapConnection == nil {
		return errors.New("Initialization error")
	}
	rc = C.ldap_simple_bind_s(tconf.ldapConnection, C.CString(tconf.login), C.CString(tconf.passwd))
	if rc != LDAP_SUCCESS {
		return errors.New("Binding error: " + C.GoString(C.ldap_err2string(rc)))
	}
	c := 0
	for c < len(mails) {
		rc := C.ldap_delete_s(tconf.ldapConnection, C.CString(mails[c]))
		if rc != 0 {
			return errors.New("Deletion error: " + C.GoString(C.ldap_err2string(rc)))
		}
		c++
	}
	return nil
}

func initialize(conf *ldap_conf) error {
	if setOptions(nil) != nil {
		return errors.New("Options error")
	}
	rc := C.ldap_initialize(&conf.ldapConnection, C.CString(conf.host))
	if conf.ldapConnection == nil || rc != LDAP_SUCCESS {
		return errors.New("Initialization error: " + C.GoString(C.ldap_err2string(rc)))
	}
	rc = C.ldap_simple_bind_s(conf.ldapConnection, C.CString(conf.login), C.CString(conf.passwd))
	if rc != LDAP_SUCCESS {
		return errors.New("Binding error: " + C.GoString(C.ldap_err2string(rc)))
	}
	return nil
}

// Checks if there is at least one sam account available, to use it to create a new user instead of generating a new sam account
func checkSamAvailability(ldapConnection *ldap.Conn) (error, string, int) {
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
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

func createNewUser(conf2 ldap_conf, params AccountParams, count int, mods [3]*C.LDAPModStr, ldapConnection *ldap.Conn) (*nano.Response, error) {
	var sam string
	if !test_password(params.Password) {
		return nano.JSONResponse(400, hash{
			"error": "Password does not meet minimum requirements",
		}), nil
	}

	dn := "cn=" + fmt.Sprintf("%d", count+1) + "," + conf2.ou

	rc := C._ldap_add(conf2.ldapConnection, C.CString(dn), &mods[0])

	if rc != LDAP_SUCCESS {
		return nil, errors.New("Adding error: " + C.GoString(C.ldap_err2string(rc)))
	}
	pwd := encodePassword(params.Password)
	modify := ldap.NewModifyRequest("cn=" + fmt.Sprintf("%d", count+1) + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)}) // field where the windows password is stored
	modify.Replace("userAccountControl", []string{"512"})
	err := ldapConnection.Modify(modify)
	if err != nil {
		return nil, errors.New("Modify error: " + err.Error())
	}
	ldapConnection, err = ldap.DialTLS("tcp", conf.LDAPServer.Host,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, errors.New("Dial error: " + err.Error())
	}
	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+fmt.Sprintf("%d", count+1)+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return nil, errors.New("Search error: " + err.Error())
	}
	for _, entry := range sr.Entries {
		module.Log.Info(entry.GetAttributeValue("sAMAccountName"))
		sam = entry.GetAttributeValue("sAMAccountName")
	}
	return nano.JSONResponse(200, hash{
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
	modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)})
	modify.Replace("userAccountControl", []string{"512"})
	modify.Replace("mail", []string{params.UserEmail})
	err := ldapConnection.Modify(modify)
	if err != nil {
		return nil, errors.New("Modify error: " + err.Error())
	}

	ldapConnection, err = ldap.DialTLS("tcp", conf.ServerURL[8:],
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, errors.New("Dial error: " + err.Error())
	}
	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+cn+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return nil, errors.New("Search error: " + err.Error())
	}
	for _, entry := range sr.Entries {
		module.Log.Info(entry.GetAttributeValue("sAMAccountName"))
		sam = entry.GetAttributeValue("sAMAccountName")
	}
	return nano.JSONResponse(200, hash{
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
		return nil, err
	}

	if len(req.Params["user_id"]) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "user id is missing",
		}), nil
	}

	params.UserEmail = req.Params["user_id"]

	bindusername := conf.Username
	bindpassword := conf.Password
	c := 0
	for i, val := range conf.ServerURL { //Passing letters/symbols before IP adress ( ex : ldaps:// )
		if unicode.IsDigit(val) {
			c = i
			break
		}
	}
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:],
		&tls.Config{
			InsecureSkipVerify: true,
		})

	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}

	defer ldapConnection.Close()

	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+params.UserEmail+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return nil, errors.New("Search error: " + err.Error())
	}

	var cn string
	if len(sr.Entries) != 1 {
		return nil, errors.New("invalid Email")
	}
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")

	}
	pwd := encodePassword(params.Password)

	modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)})
	err = ldapConnection.Modify(modify)
	if err != nil {
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
		return nil, err
	}

	// openLDAP and CGO needed here to add a new user
	var tconf ldap_conf
	tconf.host = conf.LDAPServer.Scheme + "://" + conf.LDAPServer.Host
	tconf.login = conf.Username
	tconf.passwd = conf.Password
	tconf.ou = "OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com"
	err = initialize(&tconf)
	if err != nil {
		return nil, err
	}
	var mods [3]*C.LDAPModStr
	var modClass, modCN C.LDAPModStr
	var vclass [5]*C.char
	var vcn [4]*C.char
	modClass.mod_op = 0
	modClass.mod_type = C.CString("objectclass")
	vclass[0] = C.CString("top")
	vclass[1] = C.CString("person")
	vclass[2] = C.CString("organizationalPerson")
	vclass[3] = C.CString("User")
	vclass[4] = nil
	modClass.mod_vals = &vclass[0]

	modCN.mod_op = 0
	modCN.mod_type = C.CString("mail")
	vcn[0] = C.CString(params.UserEmail)
	vcn[1] = nil
	modCN.mod_vals = &vcn[0]

	mods[0] = &modClass
	mods[1] = &modCN
	mods[2] = nil

	bindusername := conf.Username
	bindpassword := conf.Password
	// return "", to ldap go API to set the password
	ldapConnection, err := ldap.DialTLS("tcp", conf.LDAPServer.Host,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, errors.New(">> DialTLS failed: " + err.Error())
	}
	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}

	defer ldapConnection.Close()

	err, cn, count := checkSamAvailability(ldapConnection) // if an account is disabled, this function will look for his CN
	if err != nil {
		return nil, err
	}

	// if no disabled accounts were found, a real new user is created
	if cn == "" {
		res, err := createNewUser(tconf, params, count, mods, ldapConnection)
		if err != nil {
			return nil, err
		}
		// freeing various structures needed for adding entry with OpenLDAP
		C.free(unsafe.Pointer(vclass[0]))
		C.free(unsafe.Pointer(vclass[1]))
		C.free(unsafe.Pointer(vclass[2]))
		C.free(unsafe.Pointer(vclass[3]))
		C.free(unsafe.Pointer(vcn[0]))
		C.free(unsafe.Pointer(modCN.mod_type))
		//C._ldap_mods_free(&mods[0], 1)   Should work but doesnt
		return res, nil
	}

	// if a disabled account is found, modifying this account instead of creating a new one
	return recycleSam(params, ldapConnection, cn)
}

func forcedisableAccount(req nano.Request) (*nano.Response, error) {
	userId := req.Params["user_id"]

	if len(userId) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "User id is missing",
		}), nil
	}

	bindusername := conf.Username
	bindpassword := conf.Password
	c := 0
	for i, val := range conf.ServerURL { // passing letters/symbols before IP adress ( ex : ldaps:// )
		if unicode.IsDigit(val) {
			c = i
			break
		}
	}
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:],
		&tls.Config{
			InsecureSkipVerify: true,
		})

	if err != nil {
		return nil, errors.New("DialTLS error: " + err.Error())
	}
	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return nil, errors.New("Binding error: " + err.Error())
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+userId+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return nil, errors.New("Searching error: " + err.Error())
	}

	if len(sr.Entries) != 1 {
		// means entered mail was not valid, or several user have the same mail
		return nil, errors.New("Email does not match any user, or several users have the same mail adress")
	}
	var cn string
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")
	}
	modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("userAccountControl", []string{"514"}) // 512 is a normal account, 514 is disabled ( 512 + 0x0002 )
	err = ldapConnection.Modify(modify)
	if err != nil {
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
