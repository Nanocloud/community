package users

import (
	"log"
	"testing"
)

var (
	activated = true
	email     = "foo@nanocloud.com"
	firstName = "Admin"
	lastName  = "Nanocloud"
	password  = "secret"
	isAdmin   = false
	id        = ""
)

func getUser(id string, error string) *User {
	user, err := GetUser(id)

	if err != nil {
		log.Fatalf("Cannot get the user: %v", err.Error())
	}
	if user == nil {
		log.Fatalf(error)
	}
	return user
}

func compareUser(user *User) {
	switch {
	case user.Activated != activated:
		log.Fatalf("'user.Activated' field doesn't match the inserted value")
	case user.Email != email:
		log.Fatalf("'user.Email' field doesn't match the inserted value")
	case user.FirstName != firstName:
		log.Fatalf("'user.FirstName' field doesn't match the inserted value")
	case user.LastName != lastName:
		log.Fatalf("'user.LastName' field doesn't match the inserted value")
	case user.Password != "":
		log.Fatalf("'user.Password' field should be empty")
	case user.IsAdmin != isAdmin:
		log.Fatalf("'user.IsAdmin' field doesn't match the inserted value")
	}
}

func TestCreateUser(t *testing.T) {
	user, err := CreateUser(activated, email, firstName, lastName, password, isAdmin)

	if err != nil {
		log.Fatalf("Cannot create the user: %v", err.Error())
		return
	}
	if user == nil {
		log.Fatalf("The user was not created")
		return
	}

	id = user.Id
	compareUser(user)
}

func TestGetUserFromEmailPassword(t *testing.T) {
	user, err := GetUserFromEmailPassword(email, password)

	if err != nil {
		t.Errorf("Cannot get the user from his email/password: %v", err.Error())
		return
	}
	if user == nil {
		t.Errorf("No error was returned but get a nil user")
		return
	}
	compareUser(user)
}

func TestDisableUser(t *testing.T) {
	err := DisableUser(id)

	if err != nil {
		t.Errorf("Cannot disable user: %v", err.Error())
	}

	user := getUser(id, "Nil user was returned")
	if user.Activated != false {
		t.Errorf("'user.Activated' field should be false")
	}
}

func TestDeleteUser(t *testing.T) {
	err := DeleteUser(id)

	if err != nil {
		t.Errorf("Cannot delete the user: %v", err.Error())
		return
	}

	exists, err := UserExists(id)
	if err != nil {
		t.Errorf("An error was returned by UserExists: %v", err.Error())
		return
	}
	if exists != false {
		log.Fatalf("User exists even after deletion")
	}
}
