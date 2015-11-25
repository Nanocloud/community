/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func extractTarGzFile(srcFile string) {

	var (
		name     string
		filePath string
		mode     os.FileMode = 0755
	)

	f, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("Can't open compressed archive %s, error: %v", srcFile, err)
		os.Exit(1)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("Can't open archive %s, error: %v", srcFile, err)
		os.Exit(1)
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Archive error: %v", err)
			os.Exit(1)
		}

		name = header.Name
		filePath = filepath.Join(prefix, name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(filePath, mode)
		case tar.TypeReg:
			file, err := os.Create(filePath)
			if err != nil {
				log.Print("Error unpacking file %s: %v", filePath, err)
			}
			io.Copy(file, tarReader)
			os.Chmod(filePath, mode)
			file.Close()
		case tar.TypeSymlink:
			err := os.Symlink(header.Linkname, filePath)
			if err != nil {
				log.Printf("Error creating symlink %s : %s", filePath, err)
			}
		default:
			log.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
}

func jsonRpcRequest(url string, method string, params map[string]string) {
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params": []map[string]string{
			0: params,
		},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	result := make(map[string]interface{})
	// TODO Check result and do something about it
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func unpackAsset(asset string, destination string) {
	data, err := Asset(asset)
	if err != nil {
		log.Fatalf("Unpack asser errot: %v", err)
		os.Exit(1)
	}

	ioutil.WriteFile(destination, data, 0755)
}
