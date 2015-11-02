package main

import (
	"github.com/natefinch/pie"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"strings"
)

type plugin struct {
	name   string
	client *rpc.Client
}

func main() {
	path := "plug"
	c, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Strerr, path)
	if err != nil {
		log.Println(err)
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	p := plugin{
		name:   name,
		client: c,
	}
	reply := true
	if err != nil {
		log.Println(err)
	}
}
