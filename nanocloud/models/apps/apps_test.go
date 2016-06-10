package apps

import (
	"bytes"
	"log"
	"testing"

	"fmt"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/satori/go.uuid"
)

var (
	id                    = ""
	collectionName        = "foo"
	alias                 = "notepad_plusplus_64.exe"
	displayName           = "Notepad PlusPlus"
	filePath              = "/bin/zsh"
	path                  = "/bin/zsh"
	iconContents   []byte = nil
	user, _               = users.GetUserFromEmailPassword("admin@nanocloud.com", "Nanocloud123+")
)

func getApp(id string, error string) *App {
	app, err := GetApp(id)

	if err != nil {
		log.Fatalf("Cannot get the app: %v", err.Error())
	}
	if app == nil {
		log.Fatalf(error)
	}
	return app
}

func compareApp(app *App) {
	switch {
	case app.Id == "":
		log.Fatalf("'app.Id' field is empty")
	case app.CollectionName != collectionName:
		log.Fatalf("'app.CollectionName' field doesn't match the inserted value")
	case app.Alias != alias:
		log.Fatalf("'app.Alias' field doesn't match the inserted value")
	case app.DisplayName != displayName:
		log.Fatalf("'app.DisplayName' field doesn't match the inserted value")
	case app.FilePath != filePath:
		log.Fatalf("'app.FilePath' field should be empty")
	case app.Path != "":
		log.Fatalf("'app.Path' field should be empty")
	case !bytes.Equal(app.IconContents, iconContents):
		fmt.Println("'app.IconContents' field doesn't match the inserted value")
	}
}

func TestPublishApp(t *testing.T) {
	//app := App{id, collectionName, alias, displayName, filePath, path, iconContents}

	//if user == nil {
	//	log.Fatalf("Administrator account is nil")
	//}
	//
	//err := PublishApp(user, &app)
	//if err != nil {
	//	log.Fatalf("Cannot publish the app: %v", err.Error())
	//	return
	//}
}

func TestGetApp(t *testing.T) {
	fmt.Printf("Testing GetApp()... ")
	id = uuid.NewV4().String()

	_, err := db.Query(
		`INSERT INTO apps
		(id, collection_name, alias, display_name, file_path, icon_content)
		VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::bytea)
		`,
		id, collectionName, alias, displayName, filePath, iconContents,
	)
	if err != nil {
		t.Errorf("Cannot create the application")
	}
	app := getApp(id, "Nil app was returned")
	compareApp(app)
	fmt.Printf("OK\r\n")
}

func TestChangeName(t *testing.T) {
	fmt.Printf("Testing ChangeName()... ")
	displayName = "Notepad++"
	err := ChangeName(id, displayName)
	if err != nil {
		t.Errorf("Can't update the application name: %s", err.Error())
	}
	app := getApp(id, "Nil app was returned")
	compareApp(app)
	fmt.Printf("OK\r\n")
}

func TestGetUserApps(t *testing.T) {
	fmt.Printf("Testing GetUserApps()... ")
	id = uuid.NewV4().String()

	_, err := db.Query(
		`INSERT INTO apps
		(id, collection_name, alias, display_name, file_path, icon_content)
		VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::bytea)
		`,
		id, collectionName, "Libre Office", displayName, filePath, iconContents,
	)
	if err != nil {
		t.Errorf("Cannot create the application")
	}

	apps, err := GetUserApps(user.GetID())
	if err != nil {
		t.Errorf("Unable to get user apps")
	}

	for _, app := range apps {
		if app == nil {
			t.Errorf("A nil app was returned")
		}
		id = app.Id
		collectionName = app.CollectionName
		alias = app.Alias
		displayName = app.DisplayName
		filePath = app.FilePath
		iconContents = app.IconContents
		compareApp(app)
		db.Query(`DELETE FROM apps where id=$1::varchar`, id)
	}
	fmt.Printf("OK\r\n")
}
