package apps

import (
	"log"
	"testing"

	"github.com/Nanocloud/community/nanocloud/models/users"
	uuid "github.com/satori/go.uuid"
)

var (
	app_num               = 0
	id                    = ""
	collectionName        = "foo"
	alias                 = "notepad_plusplus_64.exe"
	displayName           = "Notepad PlusPlus"
	filePath              = "/bin/zsh"
	path                  = "/bin/zsh"
	iconContents   []byte = nil
	user                  = &users.User{}
	list_apps             = []*App{}
)

func init() {
	new_user, err := users.CreateUser(
		true,
		"new_user@nanocloud.com",
		"Test",
		"user",
		"secret",
		false,
	)

	if err != nil {
		log.Panicln("Can't create new account:", err.Error())
	}
	if new_user == nil {
		log.Panicln("Can't create new account")
	}
	user = new_user
}

func getApp(app_id string, error string) *App {
	get_app, err := GetApp(app_id)

	if err != nil {
		log.Panicln("Cannot get the app:", err.Error())
	}
	if get_app == nil {
		log.Panicf(error)
	}

	list_apps = append(list_apps, get_app)
	id = get_app.Id
	return get_app
}

func compareApp(get_app *App, i int) {
	switch {
	case get_app.Id == "":
		log.Fatalln("'app.Id' field is empty")
	case get_app.CollectionName != list_apps[i].CollectionName:
		log.Fatalln("'app.CollectionName' field doesn't match the inserted value")
	case get_app.Alias != list_apps[i].Alias:
		log.Fatalln("'app.Alias' field doesn't match the inserted value")
	case get_app.DisplayName != list_apps[i].DisplayName:
		log.Fatalln("'app.DisplayName' field doesn't match the inserted value")
	case get_app.FilePath != list_apps[i].FilePath:
		log.Fatalln("'app.FilePath' field should be empty")
	}
}

func TestCreateApp(t *testing.T) {
	new_app := &App{id, collectionName, alias, displayName, filePath, path, iconContents}

	new_app, err := CreateApp(new_app)
	if err != nil {
		log.Fatalln("Cannot create the app:", err.Error())
	}

	new_app = getApp(new_app.GetID(), "Can't get the created application")
	compareApp(new_app, app_num)
	app_num++
}

func TestGetApp(t *testing.T) {
	alias = "LibreOffice.exe"
	displayName = "Libre Office"
	new_app := &App{Id: "", CollectionName: collectionName, Alias: alias, DisplayName: displayName, FilePath: filePath}

	new_app, err := CreateApp(new_app)
	if err != nil {
		log.Fatalln("Cannot create the app:", err.Error())
	}
	_ = getApp(new_app.GetID(), "Nil app was returned")
	compareApp(new_app, app_num)
	app_num++
}

func TestChangeName(t *testing.T) {
	displayName = "Notepad++"
	err := ChangeName(id, displayName)
	if err != nil {
		t.Errorf("Can't update the application name: %s", err.Error())
	}

	get_app, err := GetApp(id)
	if err != nil {
		t.Errorf("Can't get the updated app: %s", err.Error())
	}
	if get_app == nil {
		t.Error("Nil app was returned")
	}
	list_apps[app_num-1].DisplayName = displayName
	compareApp(get_app, app_num-1)
}

func TestGetUserApps(t *testing.T) {
	id = uuid.NewV4().String()
	new_app := &App{Id: id, CollectionName: collectionName, Alias: "Libre Office", DisplayName: displayName, FilePath: filePath}
	i := 0

	new_app, err := CreateApp(new_app)
	if err != nil {
		log.Fatalln("Cannot create the app:", err.Error())
	}
	list_apps = append(list_apps, new_app)
	app_num++

	apps, err := GetUserApps(user.GetID())
	if err != nil {
		t.Error("Unable to get user apps")
	}

	for _, get_app := range apps {
		if get_app == nil {
			t.Error("A nil app was returned")
		}
		if get_app.Alias != "hapticDesktop" {
			compareApp(get_app, i)
			err = get_app.Delete()
			if err != nil {
				log.Fatalln("Can't delete application:", err.Error())
			}
			i++
		}
	}

	err = users.DeleteUser(user.GetID())
	if err != nil {
		t.Errorf("Can't delete user: %s\n", err.Error())
	}
}
