package users

import (
	"fmt"
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
	fmt.Printf("Testing CreateUser()... ")
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
	fmt.Printf("OK\r\n")
}

func TestGetUserFromEmailPassword(t *testing.T) {
	fmt.Printf("Testing GetUserFromEmailPassword()... ")
	user, err := GetUserFromEmailPassword(email, password)

	if err != nil {
		t.Errorf("Cannot get the user from his email/password: %v", err.Error())
		return
	}
	if user == nil {
		t.Errorf("No error was returned but get a nil user\r\n")
		return
	}
	compareUser(user)
	fmt.Printf("OK\r\n")
}

func TestUpdateUserEmail(t *testing.T) {
	fmt.Printf("Testing UpdateUserEmail()... ")
	err := UpdateUserEmail(id, "bar@nanocloud.com")

	if err != nil {
		t.Errorf("Cannot update user email: %v", err.Error())
	}

	email = "bar@nanocloud.com"
	user := getUser(id, "Nil user was returned")
	compareUser(user)
	fmt.Printf("OK\r\n")
}

func TestUpdateUserFirstName(t *testing.T) {
	fmt.Printf("Testing UpdateUserFirstName()... ")
	err := UpdateUserFirstName(id, "Nano")

	if err != nil {
		t.Errorf("Cannot update user first name: %v", err.Error())
	}

	firstName = "Nano"
	user := getUser(id, "Nil user was returned")
	compareUser(user)
	fmt.Printf("OK\r\n")
}

func TestUpdateUserLastName(t *testing.T) {
	fmt.Printf("Testing UpdateUserLastName()... ")
	err := UpdateUserLastName(id, "Cloud")

	if err != nil {
		t.Errorf("Cannot update user last name: %v", err.Error())
	}

	lastName = "Cloud"
	user := getUser(id, "Nil user was returned")
	compareUser(user)
	fmt.Printf("OK\r\n")
}

func TestUpdateUserPassword(t *testing.T) {
	fmt.Printf("Testing UpdateUserPassword()... ")
	err := UpdateUserPassword(id, "foobar")
	if err != nil {
		t.Errorf("Cannot update user password %v", err.Error())
	}

	password = "foobar"
	user := getUser(id, "Nil user was returned")
	compareUser(user)
	fmt.Printf("OK\r\n")
}

func TestDisableUser(t *testing.T) {
	fmt.Printf("Testing DisableUser()... ")
	err := DisableUser(id)

	if err != nil {
		t.Errorf("Cannot disable user: %v", err.Error())
	}

	user := getUser(id, "Nil user was returned")
	if user.Activated != false {
		t.Errorf("'user.Activated' field should be false\r\n")
	}
	fmt.Printf("OK\r\n")
}

func TestDeleteUser(t *testing.T) {
	fmt.Printf("Testing DeleteUser()... ")
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
	fmt.Printf("OK\r\n")
}
