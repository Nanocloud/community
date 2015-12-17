package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/natefinch/pie"
	"gopkg.in/fsnotify.v1"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	moduleName = "nanocloud"
)

var (
	plugins = make(map[string]plugin)
)

func setupDb() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	// oauth_clients table
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'oauth_clients'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("[nanocloud] oauth_clients table already set up\n")
	} else {
		rows, err = db.Query(
			`CREATE TABLE oauth_clients (
				id      serial PRIMARY KEY,
				name    varchar(255) UNIQUE,
				key     varchar(255) UNIQUE,
				secret  varchar(255)
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create oauth_clients table: %s\n", err)
			return err
		}
		defer rows.Close()

		rows, err = db.Query(
			`INSERT INTO oauth_clients
			(name, key, secret)
			VALUES (
				'Nanocloud',
				'9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae',
				'9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341'
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create default oauth_clients: %s\n", err)
			return err
		}
		defer rows.Close()
	}

	// oauth_access_tokens table
	rows, err = db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'oauth_access_tokens'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("[nanocloud] oauth_access_tokens table already set up\n")
	} else {
		rows, err = db.Query(
			`CREATE TABLE oauth_access_tokens (
				id                serial PRIMARY KEY,
				token             varchar(255) UNIQUE,
				oauth_client_id   integer REFERENCES oauth_clients (id),
				user_id           varchar(255)
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create oauth_access_tokens table: %s\n", err)
			return err
		}
		defer rows.Close()
	}
	return nil
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	initConf()
	setupDb()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Unable to run the watcher. ", err)
	}
	defer w.Close()

	runningPlugins := launchExistingPlugins()
	go watchPlugins(w, runningPlugins)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", conf.FrontDir)
	e.Get("/api/me", getMeHandler)
	e.Get("/api/version", getVersionHandler)
	e.Any("/api/*", genericHandler)
	e.Any("/oauth/*", genericHandler)

	addr := ":" + conf.Port
	log.Info("Server running at ", addr)
	e.Run(addr)
}

// get plugins list from the running directory
func getPlugins() []string {
	dir, err := os.Open(conf.RunDir[:len(conf.RunDir)-1])
	if err != nil {
		log.Error("Unable to open the running folder. ", err)
		return nil
	}
	defer dir.Close()
	var filenames []string
	fis, err := dir.Readdir(-1) // -1 means return all the FileInfos
	if err != nil {
		log.Error("Unable to get filenames of the running folder. ", err)
		return nil
	}
	for _, fileinfo := range fis {
		if !fileinfo.IsDir() && !strings.HasSuffix(fileinfo.Name(), ".tar.gz") {
			filenames = append(filenames, fileinfo.Name())
		}
	}

	sort.Strings(filenames)

	return filenames
}

// load all available plugins
func launchExistingPlugins() []string {
	var runningPlugins []string
	plugs := getPlugins()
	for _, plugin := range plugs {
		runningPlugins = addPlugin(runningPlugins, plugin)
		copyFile(conf.RunDir+plugin, conf.InstDir+plugin)
	}
	return runningPlugins
}

// check if the plugin is currently running
func isRunning(runningPlugins []string, name string) bool {
	for _, val := range runningPlugins {
		if val == name {
			return true
		}
	}
	return false
}

// remove a plugin
func removePlugin(runningPlugins []string, name string) []string {
	for i, val := range runningPlugins {
		if val == name {
			closePlugin(conf.RunDir + name)
			runningPlugins = append(runningPlugins[:i], runningPlugins[i+1:]...)
			log.Println("deleted plugin from slice")
		}
	}
	return runningPlugins
}

// add a plugin
func addPlugin(runningPlugins []string, name string) []string {

	loadPlugin(conf.RunDir + name)
	runningPlugins = append(runningPlugins, name)

	return runningPlugins
}

// delete a plugin from disk
func deletePlugin(path string) error {
	oldpath := path
	if _, ok := plugins[path]; !ok {
		path = conf.StagDir + path[strings.LastIndex(path, "/")+1:]
	}

	if path == oldpath {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(oldpath)
		if err != nil {
			return err
		}
	}

	return nil
}

// handle the case when a plugin is dropped in the staging folder
func createEvent(runningPlugins []string, name string, fullpath string, sourcefile string) []string {

	// if another version of the plugin is already running
	if isRunning(runningPlugins, name) {
		runningPlugins = removePlugin(runningPlugins, name)
		err := loadPlugin(fullpath)
		// verify if the plugin can be loaded
		if err == nil && plugins[fullpath].Check() == true {
			runningPlugins = append(runningPlugins, name)

			deletePlugin(conf.RunDir + name)
			err := os.Rename(conf.StagDir+name, conf.RunDir+name)
			if err != nil {
				log.Println(err)
			}

			copyFile(conf.RunDir+name, conf.InstDir+name)
			deleteOldFront(sourcefile)
			unpackFront(sourcefile)
			err = os.Rename(sourcefile, conf.RunDir+sourcefile[strings.LastIndex(sourcefile, "/")+1:])

			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("New plugin encountered an error")
			err := loadPlugin(conf.RunDir + name)
			if err != nil {
				log.Println("error loading plugin")
				log.Println(err)
			}
		}

	} else {

		// ok, it's a new plugin

		// handle the binary
		err := os.Rename(conf.StagDir+name, conf.RunDir+name)
		if err != nil {
			log.Error("Unable to move the binary file to the running folder. ", err)
		}
		_, err = os.Create(conf.InstDir + name)
		if err != nil {
			log.Error("Unable to create a touch file in the installed folder. ", err)
		}
		// handle the front
		unpackFront(sourcefile)
		err = os.Rename(sourcefile, conf.RunDir+sourcefile[strings.LastIndex(sourcefile, "/")+1:])
		if err != nil {
			log.Error("Unable to move the tarball to the running folder. ", err)
		}

		runningPlugins = addPlugin(runningPlugins, name)
	}
	return runningPlugins
}

func deleteTar(name string) {
	err := os.Remove(conf.RunDir + name[strings.LastIndex(name, "/")+1:])
	if err != nil {
		log.Println(err)
	}
}

func watchPlugins(w *fsnotify.Watcher, runningPlugins []string) {
	w.Add(conf.StagDir)
	w.Add(conf.InstDir)
	for {
		select {

		case evt := <-w.Events:
			switch evt.Op {
			case fsnotify.Create:
				if evt.Name[:strings.LastIndex(evt.Name, "/")+1] == conf.StagDir {
					unpackGo(evt.Name, runningPlugins)
				}
			case fsnotify.Remove:
				deleteTar(evt.Name + ".tar.gz")
				closePlugin(conf.RunDir + evt.Name[strings.LastIndex(evt.Name, "/")+1:])
				deletePlugin(conf.RunDir + evt.Name[strings.LastIndex(evt.Name, "/")+1:])
			}

		case err := <-w.Errors:
			log.Println("watcher crashed:", err)

		}
	}
}

func loadPlugin(path string) error {

	c, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	if err != nil {
		log.Printf("Error running plugin %s: %s", path, err)
		return err
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	p := plugin{
		name:   name,
		client: c,
	}
	p.Plug()

	plugins[path] = p
	return nil
}

func closePlugin(path string) {
	if _, ok := plugins[path]; !ok {
		path = conf.StagDir + path[strings.LastIndex(path, "/")+1:]
		if _, ok := plugins[path]; !ok {
			log.Println("Plugin not found for deletion")
			return
		}
	}

	plugins[path].Unplug()

	delete(plugins, path)
}
