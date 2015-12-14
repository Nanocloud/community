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
	"gopkg.in/ldap.v2"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	"github.com/natefinch/pie"
	"github.com/streadway/amqp"
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

var (
	name = "ldap"
	srv  pie.Server
)

type api struct{}

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

// Structure used for exchanges with the core, faking a request/responsewriter
type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	Status   int
	HeadVals map[string]string
}

// Strucutre used in return messages sent to RabbitMQ
type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

// Setting LDAP version and referrals
func setOptions(ldapConnection *C.LDAP) error {
	var version C.int
	var opt C.int
	version = LDAP_VERSION3
	opt = 0
	err := C.ldap_set_option(ldapConnection, LDAP_OPT_PROTOCOL_VERSION, unsafe.Pointer(&version))
	if err != LDAP_SUCCESS {
		return answerWithError("Options settings error: "+C.GoString(C.ldap_err2string(err)), nil)
	}
	err = C.ldap_set_option(ldapConnection, LDAP_OPT_REFERRALS, unsafe.Pointer(&opt))
	if err != LDAP_SUCCESS {
		return answerWithError("Options settings error: "+C.GoString(C.ldap_err2string(err)), nil)
	}
	return nil

}

func answerWithError(msg string, e error) error {

	answer := "plugin Ldap: " + msg
	if e != nil {
		answer += e.Error()
	}

	log.Error(answer)

	return errors.New(answer)
}

func listUsers(args PlugRequest, reply *PlugRequest, id string) error {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"
	reply.Status = 500
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[8:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return answerWithError("Dial error: ", err)
	}
	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return answerWithError("Binding error: ", err)
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
		return answerWithError("Search error: ", err)
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
	g, _ := json.Marshal(res)
	reply.Body = string(g)
	reply.Status = 200
	return nil
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
		return answerWithError("Options settings error: "+C.GoString(C.ldap_err2string(err)), nil)
	}

	err = C.ldap_set_option(tconf.ldapConnection, LDAP_OPT_REFERRALS, unsafe.Pointer(&v))
	if err != LDAP_SUCCESS {
		return answerWithError("Deletion error: "+C.GoString(C.ldap_err2string(err)), nil)
	}

	rc := C.ldap_initialize(&tconf.ldapConnection, C.CString(tconf.host+":636"))
	if tconf.ldapConnection == nil {
		return answerWithError("Initialization error: ", nil)
	}
	rc = C.ldap_simple_bind_s(tconf.ldapConnection, C.CString(tconf.login), C.CString(tconf.passwd))
	if rc != LDAP_SUCCESS {
		return answerWithError("Binding error: "+C.GoString(C.ldap_err2string(rc)), nil)
	}
	c := 0
	for c < len(mails) {
		rc := C.ldap_delete_s(tconf.ldapConnection, C.CString(mails[c]))
		if rc != 0 {
			return answerWithError("Deletion error: "+C.GoString(C.ldap_err2string(rc)), nil)
		}
		c++
	}
	return nil

}
func initialize(conf *ldap_conf) error {

	if setOptions(nil) != nil {
		return answerWithError("Options error", nil)
	}
	rc := C.ldap_initialize(&conf.ldapConnection, C.CString(conf.host+":636"))
	if conf.ldapConnection == nil {
		return answerWithError("Initialization error: "+C.GoString(C.ldap_err2string(rc)), nil)
	}
	rc = C.ldap_simple_bind_s(conf.ldapConnection, C.CString(conf.login), C.CString(conf.passwd))
	if rc != LDAP_SUCCESS {
		return answerWithError("Binding error: "+C.GoString(C.ldap_err2string(rc)), nil)
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
		return answerWithError("Search error: ", err), "", 0
	}
	count := len(sr.Entries)
	cn := ""
	for _, entry := range sr.Entries {
		h, err := strconv.Atoi(entry.GetAttributeValue("userAccountControl"))
		if err != nil {
			return answerWithError("Atoi conversion error: ", err), "", 0
		}
		if h&0x0002 == 0 { //0x0002 means disabled account
		} else {
			cn = entry.GetAttributeValue("cn")
			break
		}
	}
	return nil, cn, count
}

func createNewUser(conf2 ldap_conf, params AccountParams, count int, mods [3]*C.LDAPModStr, ldapConnection *ldap.Conn, reply *PlugRequest) error {

	if !test_password(params.Password) {
		reply.Status = 400
		return answerWithError("Password does not meet minimum requirements", nil)

	}
	dn := "cn=" + fmt.Sprintf("%d", count+1) + "," + conf2.ou

	rc := C._ldap_add(conf2.ldapConnection, C.CString(dn), &mods[0])

	if rc != LDAP_SUCCESS {
		return answerWithError("Adding error: "+C.GoString(C.ldap_err2string(rc)), nil)
	}
	pwd := encodePassword(params.Password)
	modify := ldap.NewModifyRequest("cn=" + fmt.Sprintf("%d", count+1) + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)}) // field where the windows password is stored
	modify.Replace("userAccountControl", []string{"512"})
	err := ldapConnection.Modify(modify)
	if err != nil {
		return answerWithError("Modify error: ", err)
	}
	ldapConnection, err = ldap.DialTLS("tcp", conf.ServerURL[8:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return answerWithError("Dial error: ", err)
	}
	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return answerWithError("Binding error: ", err)
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
		return answerWithError("Search error: ", err)
	}
	for _, entry := range sr.Entries {
		log.Info(entry.GetAttributeValue("sAMAccountName"))
	}
	return nil

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
func recycleSam(params AccountParams, ldapConnection *ldap.Conn, cn string) error {

	pwd := encodePassword(params.Password)
	modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)})
	modify.Replace("userAccountControl", []string{"512"})
	modify.Replace("mail", []string{params.UserEmail})
	err := ldapConnection.Modify(modify)
	if err != nil {
		return answerWithError("Modify error: ", err)
	}

	ldapConnection, err = ldap.DialTLS("tcp", conf.ServerURL[8:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return answerWithError("Dial error: ", err)
	}
	err = ldapConnection.Bind(conf.Username, conf.Password)
	if err != nil {
		return answerWithError("Binding error: ", err)
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
		return answerWithError("Search error: ", err)
	}
	for _, entry := range sr.Entries {
		log.Info(entry.GetAttributeValue("sAMAccountName"))
	}
	return nil
}

func modifyPassword(args PlugRequest, reply *PlugRequest, id string) error {

	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	reply.Status = 500

	var params AccountParams

	if e := json.Unmarshal([]byte(args.Body), &params); e != nil {
		reply.Status = 400
		return answerWithError("modify password failed: ", e)
	}
	if id != "" {
		params.UserEmail = id
	}
	bindusername := conf.Username
	bindpassword := conf.Password
	c := 0
	for i, val := range conf.ServerURL { //Passing letters/symbols before IP adress ( ex : ldaps:// )
		if unicode.IsDigit(val) {
			c = i
			break
		}
	}
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})

	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return answerWithError("Binding error: ", err)
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
		return answerWithError("Search error: ", err)
	}

	var cn string
	if len(sr.Entries) != 1 {
		return answerWithError("invalid Email", nil)
	}
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")

	}
	pwd := encodePassword(params.Password)

	modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
	modify.Replace("unicodePwd", []string{string(pwd)})
	err = ldapConnection.Modify(modify)
	if err != nil {
		return answerWithError("Password modification failed: ", err)
	}
	reply.Status = 202
	return nil

}

func addUser(args PlugRequest, reply *PlugRequest, id string) error {
	reply.Status = 201
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	var params AccountParams

	if e := json.Unmarshal([]byte(args.Body), &params); e != nil {
		reply.Status = 400
		return answerWithError("addUser() failed: ", e)
	}
	// openLDAP and CGO needed here to add a new user
	var tconf ldap_conf
	tconf.host = conf.ServerURL
	tconf.login = conf.Username
	tconf.passwd = conf.Password
	tconf.ou = "OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com"
	err := initialize(&tconf)
	if err != nil {
		return err
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
	// return to ldap go API to set the password
	c := 0
	for i, val := range conf.ServerURL { // passing letters/symbols before IP adress ( ex : ldaps:// )
		if unicode.IsDigit(val) {
			c = i
			break
		}
	}
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return answerWithError("DialTLS failed: ", err)
	}
	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return answerWithError("Binding error: ", err)
	}

	defer ldapConnection.Close()

	err, cn, count := checkSamAvailability(ldapConnection) // if an account is disabled, this function will look for his CN
	if err != nil {
		return err
	}

	// if no disabled accounts were found, a real new user is created
	if cn == "" {
		err = createNewUser(tconf, params, count, mods, ldapConnection, reply)
		if err != nil {
			return err
		}
		// freeing various structures needed for adding entry with OpenLDAP
		C.free(unsafe.Pointer(vclass[0]))
		C.free(unsafe.Pointer(vclass[1]))
		C.free(unsafe.Pointer(vclass[2]))
		C.free(unsafe.Pointer(vclass[3]))
		C.free(unsafe.Pointer(vcn[0]))
		C.free(unsafe.Pointer(modCN.mod_type))
		//C._ldap_mods_free(&mods[0], 1)   Should work but doesnt
	} else {
		// if a disabled account is found, modifying this account instead of creating a new one
		err = recycleSam(params, ldapConnection, cn)
		if err != nil {
			return err
		}

	}
	reply.Status = 201
	return nil
}

func forcedisableAccount(args PlugRequest, reply *PlugRequest, id string) error {
	reply.Status = 202
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	var params AccountParams
	if id != "" {
		params.UserEmail = id
	} else if e := json.Unmarshal([]byte(args.Body), &params); e != nil {
		reply.Status = 400
		return answerWithError("forcedisableAccount() failed ", e)
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
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})

	if err != nil {
		return answerWithError("DialTLS error: ", err)
	}
	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return answerWithError("Binding error: ", err)
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+params.UserEmail+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return answerWithError("Searching error: ", err)
	}

	if len(sr.Entries) != 1 {
		// means entered mail was not valid, or several user have the same mail
		return answerWithError("Email does not match any user, or several users have the same mail adress", nil)
	} else {
		var cn string
		for _, entry := range sr.Entries {
			cn = entry.GetAttributeValue("cn")
		}
		modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
		modify.Replace("userAccountControl", []string{"514"}) // 512 is a normal account, 514 is disabled ( 512 + 0x0002 )
		err = ldapConnection.Modify(modify)
		if err != nil {
			return answerWithError("Modify error: ", err)
		}

	}
	return nil
}

func disableAccount(args PlugRequest, reply *PlugRequest, id string) error {
	var params AccountParams

	if err := json.Unmarshal([]byte(args.Body), &params); err != nil {
		log.Error("Failed to Unmarshall the parameters of the account to disable", err)
		return err
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
	ldapConnection, err := ldap.DialTLS("tcp", conf.ServerURL[c:]+":636",
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return answerWithError("DialTLS error: ", err)
	}
	err = ldapConnection.Bind(bindusername, bindpassword)
	if err != nil {
		return answerWithError("Binding error: ", err)
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		"OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(samaccountname="+params.UserEmail+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return answerWithError("Search error: ", err)
	}

	if len(sr.Entries) != 1 { // wrong samaccount
		return answerWithError("SAMACCOUNT does not match any user", nil)
	} else {
		var cn string
		for _, entry := range sr.Entries {
			cn = entry.GetAttributeValue("cn")
		}

		modify := ldap.NewModifyRequest("cn=" + cn + ",OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")
		modify.Replace("userAccountControl", []string{"514"})
		err = ldapConnection.Modify(modify)
		if err != nil {
			return answerWithError("Modify error: ", err)
		}

	}
	return nil
}

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string) error
}{
	{`^\/api\/ldap\/users\/{0,1}$`, "POST", addUser},
	{`^\/api\/ldap\/users\/{0,1}$`, "GET", listUsers},
	{`^\/api\/ldap\/users\/(?P<id>[^\/]+)\/{0,1}$`, "PUT", modifyPassword},
	{`^\/api\/ldap\/users\/(?P<id>[^\/]+)\/disable\/{0,1}$`, "POST", forcedisableAccount},
	//	{`^\/ldap\/users\/(?P<id>[^\/]+)\/forcedisable\/{0,1}$`, "POST", forcedisableAccount},
}

// Will receive all http requests starting by /api/history from the core and chose the correct handler function
func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	var err error
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				err = val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
			} else {
				err = val.f(args, reply, "")
			}
			if err != nil {
				if reply.Status == 201 {
					reply.Status = 500
				}
			}
		}
	}
	return nil
}

// When being contacted via a queue with RabbitMQ, send a return message informing whoever called this plugin that the operation succeeded or failed
func sendReturn(msg ReturnMsg) {
	Str, err := json.Marshal(msg)
	if err != nil {
		log.Error("Faile to Marshal the rabbitmq message to return", err)
	}
	conn, err := amqp.Dial(conf.QueueURI)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	err = ch.ExchangeDeclare(
		"users_topic", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	err = ch.Publish(
		"users_topic",    // exchange
		"owncloud.users", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        Str,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent return to users")

}

// Infinite loop, awaiting for messages from RabbitMQ
func lookForMsg() {
	conn, err := amqp.Dial(conf.QueueURI)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	defer ch.Close()
	defer conn.Close()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"users_topic", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	_, err = ch.QueueDeclare(
		"ldap", // name
		true,   // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare an queue")

	err = ch.QueueBind(
		"ldap",        // queue name
		"users.*",     // routing key
		"users_topic", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := ch.Consume(
		"ldap", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		var msg Message
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Error("Failed to Unmarshal the message received: ", err)
			}
			HandleRequest(msg)

		}
	}()

	log.Printf(" [*] Waiting for messages from Users")
	<-forever
}

func HandleError(err error, mail string, method string) {
	if err != nil {
		log.Error(err)
		sendReturn(ReturnMsg{Method: method, Err: err.Error(), Plugin: "ldap", Email: mail})
	} else {
		sendReturn(ReturnMsg{Method: method, Err: "", Plugin: "ldap", Email: mail})
	}
}

// Analyze the message received from RabbitMQ
func HandleRequest(msg Message) {
	initConf()
	var args PlugRequest
	var reply PlugRequest
	var params AccountParams
	params.UserEmail = msg.Email
	params.Password = msg.Password
	userjson, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}
	args.Body = string(userjson)
	if msg.Method == "Add" {
		err = addUser(args, &reply, "")
		HandleError(err, params.UserEmail, msg.Method)
	} else if msg.Method == "disableAccount" {
		err = forcedisableAccount(args, &reply, "")
		HandleError(err, params.UserEmail, msg.Method)
	} else if msg.Method == "ChangePassword" {
		err := modifyPassword(args, &reply, "")
		HandleError(err, params.UserEmail, msg.Method)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Error(msg, err)
	}
}

// Plug the plugin to the core
func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	go lookForMsg()
	return nil
}

// Will contain various verifications for the plugin. If the core can call the function and receives "true" in the reply, it means the plugin is functionning correctly
func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

// Unplug the plugin from the core
func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {
	initConf()

	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatal("Failed to register ", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
