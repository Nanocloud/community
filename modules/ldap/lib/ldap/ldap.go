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
	"fmt"
	"net/url"
	"strconv"
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

var ConnectAndAuthWin = errors.New("Failed to connect/authenticate to windows")
var WeakPassword = errors.New("Password does not meet the minimum requirements")
var AddError = errors.New("Couldn't add user")

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

func getNumberofUsers(ldapConnection *ldap.Conn, Ou string) (error, int) {
	searchRequest := ldap.NewSearchRequest(
		Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		log.Error("Search error:  " + err.Error())
		return errors.New("Search error: " + err.Error()), 0
	}
	count := len(sr.Entries)
	return nil, count
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

func (l *Ldap) AddUser(mail, password string) (string, error) {
	var tconf ldap_conf
	tconf.host = l.LDAPServer.Scheme + "://" + l.LDAPServer.Host
	tconf.login = l.Username
	tconf.passwd = l.Password
	tconf.ou = l.Ou
	// return "", to ldap go API to set the password

	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connection to Active Directory: " + err.Error())
		return "", ConnectAndAuthWin
	}

	defer ldapConnection.Close()

	err, count := getNumberofUsers(ldapConnection, l.Ou)
	if err != nil {
		log.Error("Error while counting users: " + err.Error())
		return "", err
	}
	var sam string
	if !test_password(password) {
		return "", WeakPassword
	}

	dn := "cn=" + fmt.Sprintf("%d", count+1) + "," + tconf.ou

	req := ldap.NewAddRequest(dn)
	req.Attribute("objectclass", []string{"top", "person", "organizationalPerson", "User"})
	req.Attribute("mail", []string{mail})
	pwd := encodePassword(password)
	req.Attribute("unicodePwd", []string{string(pwd)})
	req.Attribute("userAccountControl", []string{"512"})
	err = ldapConnection.Add(req)
	if err != nil {
		log.Error("Adding error:  " + err.Error())
		return "", AddError
	}

	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(cn="+fmt.Sprintf("%d", count+1)+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		log.Error("Search error:  " + err.Error())
		return "", err
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
		return Res{}, err
	}
	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(objectGUID=*))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)
	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		return Res{}, err
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
	return res, nil
}

func (l *Ldap) ChangePassword(mail, password string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connection to Active Directory: " + err.Error())
		return err
	}

	defer ldapConnection.Close()

	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+mail+"))",
		[]string{"dn", "cn", "mail", "sAMAccountName", "userAccountControl"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		log.Error("Searching error: " + err.Error())
		return err
	}

	var cn string
	if len(sr.Entries) != 1 {
		log.Error("Invalid Email")
		return err
	}
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")
	}
	pwd := encodePassword(password)

	modify := ldap.NewModifyRequest("cn=" + cn + "," + l.Ou)
	modify.Replace("unicodePwd", []string{string(pwd)})
	err = ldapConnection.Modify(modify)
	if err != nil {
		log.Error("Password modification failed: " + err.Error())
		return err
	}
	return nil
}

func (l *Ldap) DisableUser(mail string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connecting to Active Directory: " + err.Error())
		return ConnectAndAuthWin
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+mail+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		log.Error("Searching error: " + err.Error())
		return err
	}

	if len(sr.Entries) != 1 {
		log.Error("Id does not match any user")
		// means entered mail was not valid, or several user have the same mail
		return err
	}
	var cn string
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")
	}
	modify := ldap.NewModifyRequest("cn=" + cn + "," + l.Ou)
	modify.Replace("userAccountControl", []string{"514"}) // 512 is a normal account, 514 is disabled ( 512 + 0x0002 )
	err = ldapConnection.Modify(modify)
	if err != nil {
		log.Error("Modify  error: " + err.Error())
		return err
	}
	return nil
}

func (l *Ldap) DeleteAccount(mail string) error {
	ldapConnection, err := l.DialandBind()
	if err != nil {
		log.Error("Error while connecting to Active Directory: " + err.Error())
		return err
	}

	defer ldapConnection.Close()
	searchRequest := ldap.NewSearchRequest(
		l.Ou,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person)(mail="+mail+"))",
		[]string{"userAccountControl", "cn"},
		nil,
	)

	sr, err := ldapConnection.Search(searchRequest)
	if err != nil {
		log.Error("Searching error: " + err.Error())
		return err
	}

	if len(sr.Entries) != 1 {
		log.Error("Id does not match any user")
		// means entered mail was not valid, or several user have the same mail
		return err
	}
	var cn string
	for _, entry := range sr.Entries {
		cn = entry.GetAttributeValue("cn")
	}
	del := ldap.NewDelRequest("cn="+cn+","+l.Ou, []ldap.Control{})
	err = ldapConnection.Del(del)
	if err != nil {
		log.Error("Delete  error: " + err.Error())
		return err
	}
	return nil
}
