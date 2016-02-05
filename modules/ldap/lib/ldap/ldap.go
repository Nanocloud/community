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

package ldap

import (
	"crypto/tls"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"

	"gopkg.in/ldap.v2"

	log "github.com/Sirupsen/logrus"
)

type Ldap struct {
	Username   string
	Password   string
	ServerURL  string
	Ou         string
	LDAPServer url.URL
}

var WeakPassword = errors.New("Password does not meet the minimum requirements")
var AddError = errors.New("Couldn't add user")
var AlreadyExists = errors.New("Entry already exists")
var GetUsersFailed = errors.New("Failed to retrieve users")
var ChangePwdFailed = errors.New("Failed to change the password")
var UnknownUser = errors.New("Unknown user")
var DisableFailed = errors.New("Failed to disable user")
var DeleteFailed = errors.New("Failed to delete user")

type ldap_conf struct {
	host   string
	login  string
	passwd string
	ou     string
}

type Res struct {
	Count int
	Users []map[string]string
}

func New(Username, Password, ServerURL, Ou string, LDAPServer url.URL) *Ldap {
	return &Ldap{
		Username:   Username,
		Password:   Password,
		ServerURL:  ServerURL,
		Ou:         Ou,
		LDAPServer: LDAPServer,
	}
}

func (l *Ldap) DialandBind() (*ldap.Conn, error) {
	ldapConnection, err := ldap.DialTLS("tcp", l.LDAPServer.Host,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return nil, err
	}

	err = ldapConnection.Bind(l.Username, l.Password)
	if err != nil {
		return nil, err
	}
	return ldapConnection, nil
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

func (l *Ldap) AddUser(id, password string) (string, error) {

	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connection to Active Directory: " + err.Error())
		return "", AddError
	}

	defer ldapConnection.Close()

	var sam string
	if !test_password(password) {
		return "", WeakPassword
	}

	dn := "cn=" + id + "," + l.Ou

	req := ldap.NewAddRequest(dn)
	req.Attribute("objectclass", []string{"top", "person", "organizationalPerson", "User"})
	//	req.Attribute("FirstName", []string{firstname}) FOR FUTURE USE: Add a first name and a lastname for all users
	pwd := encodePassword(password)
	req.Attribute("unicodePwd", []string{string(pwd)})
	req.Attribute("userAccountControl", []string{"512"})
	err = ldapConnection.Add(req)
	if err != nil {
		log.Error("Adding error:  " + err.Error())
		if strings.Contains(err.Error(), "Already Exists") {
			return "", AlreadyExists
		}
		return "", AddError
	}

	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+id+"))",
		[]string{"dn", "cn", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return "", AddError
	}
	for _, entry := range sr.Entries {
		log.Info(entry.GetAttributeValue("sAMAccountName"))
		sam = entry.GetAttributeValue("sAMAccountName")
	}
	return sam, nil
}

func (l *Ldap) GetUsers() (Res, error) {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connection to Active Directory: " + err.Error())
		return Res{}, GetUsersFailed
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return Res{}, GetUsersFailed
	}
	// Struct needed for JSON encoding
	var res Res
	res.Count = len(sr.Entries)
	res.Users = make([]map[string]string, res.Count)
	i := 0

	for _, entry := range sr.Entries {
		res.Users[i] = make(map[string]string, 6)
		res.Users[i]["dn"] = entry.DN
		res.Users[i]["cn"] = entry.GetAttributeValue("cn")
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
	return res, nil
}

func (l *Ldap) ChangePassword(id, password string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connection to Active Directory: " + err.Error())
		return ChangePwdFailed
	}

	defer ldapConnection.Close()
	pwd := encodePassword(password)
	modify := ldap.NewModifyRequest("cn=" + id + "," + l.Ou)
	modify.Replace("unicodePwd", []string{string(pwd)})
	err = ldapConnection.Modify(modify)
	if err != nil {
		log.Error("Password modification failed: " + err.Error())
		return ChangePwdFailed
	}
	return nil
}

func (l *Ldap) DisableUser(id string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connecting to Active Directory: " + err.Error())
		return DisableFailed
	}

	defer ldapConnection.Close()
	modify := ldap.NewModifyRequest("cn=" + id + "," + l.Ou)
	modify.Replace("userAccountControl", []string{"514"}) // 512 is a normal account, 514 is disabled ( 512 + 0x0002 )
	err = ldapConnection.Modify(modify)
	if err != nil {
		log.Error("Modify  error: " + err.Error())
		if strings.Contains(err.Error(), "No Such Object") {
			return UnknownUser
		}
		return DisableFailed
	}
	return nil
}

func (l *Ldap) DeleteAccount(id string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connecting to Active Directory: " + err.Error())
		return DeleteFailed
	}
	defer ldapConnection.Close()
	del := ldap.NewDelRequest("cn="+id+","+l.Ou, []ldap.Control{})
	err = ldapConnection.Del(del)
	if err != nil {
		log.Error("Delete  error: " + err.Error())
		if strings.Contains(err.Error(), "No Such Object") {
			return UnknownUser
		}
		return DeleteFailed
	}
	return nil
}
