package main

import (
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/pie"
	"gopkg.in/fsnotify.v1"
)

type plugin struct {
	name   string
	client *rpc.Client
}

var (
	plugins         = make(map[string]plugin)
	running_plugins []string
)

func get_plugins() []string {
	//dir, err := os.Open("plugins/running")
	dir, err := os.Open(conf.RunDir[:len(conf.RunDir)-1])
	checkErr(err)
	defer dir.Close()
	fi, err := dir.Stat()
	checkErr(err)
	filenames := make([]string, 0)
	if fi.IsDir() {
		fis, err := dir.Readdir(-1) // -1 means return all the FileInfos
		checkErr(err)
		for _, fileinfo := range fis {
			if !fileinfo.IsDir() {
				filenames = append(filenames, fileinfo.Name())
			}
		}
	}
	//log.Println("Files: ", filenames)
	return filenames
}

func checkErr(err error) {
	if err != nil {
		log.Println("Error :")
		log.Println(err)

	}
}

func launch_existing_plugins(running_plugins []string) []string {
	plugs := get_plugins()
	for _, plugin := range plugs {
		running_plugins = LoadPlugin(running_plugins, plugin)
		CopyFile(conf.RunDir+plugin, conf.InstDir+plugin)
	}
	return running_plugins
}

func main() {
	initConf()
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	router := NewRouter()
	running_plugins = make([]string, 0)
	running_plugins = launch_existing_plugins(running_plugins)

	go watchPlugins(w, running_plugins)

	log.Fatal(http.ListenAndServe(":"+conf.Port, router))

}

func ClosePlugin(running_plugins []string, name string) []string {
	for i, val := range running_plugins {
		if val == name {
			closePlugin(conf.RunDir + name)
			running_plugins = append(running_plugins[:i], running_plugins[i+1:]...)
			log.Println("deleted plugin from slice")
		}
	}
	return running_plugins
}

func LoadPlugin(running_plugins []string, name string) []string {

	loadPlugin(conf.RunDir + name)
	running_plugins = append(running_plugins, name)

	return running_plugins
}

func IsRunning(running_plugins []string, name string) bool {
	check := false
	for _, val := range running_plugins {
		if val == name {
			check = true
		}
	}
	return check
}

func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return
}

func DeletePlugin(path string) {
	oldpath := path
	if _, ok := plugins[path]; !ok {
		//		path = path[:8] + "staging" + path[15:]
		path = conf.StagDir + path[strings.LastIndex(path, "/")+1:]

	}

	if path == oldpath {
		err := os.Remove(path)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := os.Remove(oldpath)
		if err != nil {
			log.Println(err)
		}
	}
}

func CreateEvent(running_plugins []string, name string, fullpath string) []string {
	time.Sleep(1000 * time.Millisecond) // TODO, Delete that and see when file is fully copied
	if IsRunning(running_plugins, name) {
		running_plugins = ClosePlugin(running_plugins, name)
		err := loadPlugin(fullpath)
		if err != nil {
			log.Println("error loading plugin")
			log.Println(err)
		}
		if plugins[fullpath].Check() == true {
			running_plugins = append(running_plugins, name)
			log.Println("added previously existing plugin to slice")
			DeletePlugin(conf.RunDir + name)
			os.Rename(conf.StagDir+name, conf.RunDir+name)
			CopyFile(conf.RunDir+name, conf.InstDir+name) // TODO, Replace this by a simple "touch"
		} else {
			log.Println("New plugin encountered an error")
			err := loadPlugin(conf.RunDir + name)
			if err != nil {
				log.Println("error loading plugin")
				log.Println(err)
			}
		}

	} else {
		os.Rename(conf.StagDir+name, conf.RunDir+name)
		CopyFile(conf.RunDir+name, conf.InstDir+name) //TODO replace this by touch
		running_plugins = LoadPlugin(running_plugins, name)
	}
	return running_plugins
}

func watchPlugins(w *fsnotify.Watcher, running_plugins []string) {
	w.Add(conf.StagDir)
	w.Add(conf.InstDir)
	for {
		select {
		case evt := <-w.Events:
			//log.Println("fsnotify:", evt)
			switch evt.Op {
			case fsnotify.Create:
				if evt.Name[:strings.LastIndex(evt.Name, "/")+1] == conf.StagDir {
					//if evt.Name[:16] == conf.StagDir {
					//running_plugins = CreateEvent(running_plugins, evt.Name[16:], evt.Name)
					running_plugins = CreateEvent(running_plugins, evt.Name[strings.LastIndex(evt.Name, "/")+1:], evt.Name)
				}
			case fsnotify.Remove:
				//closePlugin(conf.RunDir + evt.Name[18:])
				closePlugin(conf.RunDir + evt.Name[strings.LastIndex(evt.Name, "/")+1:])
				//DeletePlugin(conf.RunDir + evt.Name[18:])
				DeletePlugin(conf.RunDir + evt.Name[strings.LastIndex(evt.Name, "/")+1:])
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

	//add_routes(name, router)

	return nil
}

func closePlugin(path string) {
	if _, ok := plugins[path]; !ok {
		//path = path[:8] + "staging" + path[15:]
		path = conf.StagDir + path[strings.LastIndex(path, "/")+1:]
		if _, ok := plugins[path]; !ok {
			log.Println("Plugin not found for deletion")
			return
		}
	}

	plugins[path].Unplug()

	delete(plugins, path)
}

func (p plugin) Plug() {
	var reply bool
	err := p.client.Call(p.name+".Plug", nil, &reply)
	if err != nil {
		log.Println("Error while calling Plug:", err)
	}
	log.Println(p.name + " plugged")
}

func (p plugin) Check() bool {
	reply := false
	err := p.client.Call(p.name+".Check", nil, &reply)
	if err != nil {
		log.Println("Error while calling Check:", err)
	}
	log.Println(p.name + " checked")
	return reply
}

func (p plugin) Unplug() {
	var reply bool
	err := p.client.Call(p.name+".Unplug", nil, &reply)
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Println("Error while calling Unplug:", err)
	}
	p.client.Close()
	log.Println(p.name + " unplugged")
}
