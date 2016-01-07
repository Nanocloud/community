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
	"github.com/Nanocloud/nano"
	"github.com/Nanocloud/oauth"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
)

// documentation for flowjs: https://github.com/flowjs/flow.js

// checkUploadHandler checks a chunk.
// If it doesn't exist then flowjs tries to upload it via uploadHandler.
func checkUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := oauth.GetUserOrFail(w, r)
	if user == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	chunkPath := filepath.Join(conf.UploadDir, user.(*nano.User).Id, "incomplete", r.FormValue("flowFilename"), r.FormValue("flowChunkNumber"))
	if _, err := os.Stat(chunkPath); err != nil {
		http.Error(w, "chunk not found", http.StatusSeeOther)
		return
	}
}

// uploadHandler tries to get and save a chunk.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	user := oauth.GetUserOrFail(w, r)
	if user == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	userPath := filepath.Join(conf.UploadDir, user.(*nano.User).Id)

	// get the multipart data
	err := r.ParseMultipartForm(2 * 1024 * 1024) // chunkSize
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chunkNum, err := strconv.Atoi(r.FormValue("flowChunkNumber"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalChunks, err := strconv.Atoi(r.FormValue("flowTotalChunks"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := r.FormValue("flowFilename")
	// module := r.FormValue("module")

	err = writeChunk(filepath.Join(userPath, "incomplete", filename), strconv.Itoa(chunkNum), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// it's done if it's not the last chunk
	if chunkNum < totalChunks {
		return
	}

	upPath := filepath.Join(userPath, filename)

	// now finish the job
	err = assembleUpload(userPath, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		module.Log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to assemble the uploaded chunks")
		return
	}
	module.Log.WithFields(log.Fields{
		"path": upPath,
	}).Info("file uploaded")

	syncOut, err := syncUploadedFile(upPath)
	if err != nil {
		module.Log.WithFields(log.Fields{
			"output": syncOut,
			"error":  err,
		}).Error("unable to scp the uploaded file to Windows")
	}
	module.Log.WithFields(log.Fields{
		"path":   upPath,
		"output": syncOut,
	}).Info("file synced")
}

func writeChunk(path, chunkNum string, r *http.Request) error {
	// prepare the chunk folder
	err := os.MkdirAll(path, 02750)
	if err != nil {
		return err
	}
	// write the chunk
	fileIn, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer fileIn.Close()
	fileOut, err := os.Create(filepath.Join(path, chunkNum))
	if err != nil {
		return err
	}
	defer fileOut.Close()
	_, err = io.Copy(fileOut, fileIn)
	return err
}

func assembleUpload(path, filename string) error {

	// create final file to write to
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	chunkDirPath := filepath.Join(path, "incomplete", filename)
	fileInfos, err := ioutil.ReadDir(chunkDirPath)
	if err != nil {
		return err
	}
	sort.Sort(byChunk(fileInfos))
	for _, fs := range fileInfos {
		src, err := os.Open(filepath.Join(chunkDirPath, fs.Name()))
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		src.Close()
		if err != nil {
			return err
		}
	}
	os.RemoveAll(chunkDirPath)

	return nil
}

func syncUploadedFile(path string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	cmd := exec.Command(filepath.Join(dir, "scripts", "copy.sh"), filepath.Join(dir, path))
	cmd.Dir = filepath.Join(dir, "scripts")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

type byChunk []os.FileInfo

func (a byChunk) Len() int      { return len(a) }
func (a byChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}
