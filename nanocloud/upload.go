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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/nanocloud/oauth"
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

	chunkPath := filepath.Join(conf.UploadDir, user.(*UserInfo).Id, "incomplete", r.FormValue("flowFilename"), r.FormValue("flowChunkNumber"))
	if _, err := os.Stat(chunkPath); err != nil {
		http.Error(w, "chunk not found", http.StatusSeeOther)
		return
	}
}

// uploadHandler tries to upload a chunk.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	user := oauth.GetUserOrFail(w, r)
	if user == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	userPath := filepath.Join(conf.UploadDir, user.(*UserInfo).Id)

	err := r.ParseMultipartForm(2 * 1024 * 1024) // chunkSize
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chunkNum := r.FormValue("flowChunkNumber")
	totalChunks := r.FormValue("flowTotalChunks")
	filename := r.FormValue("flowFilename")

	chunkDirPath := filepath.Join(userPath, "incomplete", filename)
	err = os.MkdirAll(chunkDirPath, 02750)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileIn, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileIn.Close()
	fileOut, err := os.Create(filepath.Join(chunkDirPath, chunkNum))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(fileOut, fileIn)
	fileOut.Close()

	if chunkNum == totalChunks {
		err = assembleUpload(userPath, filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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

type byChunk []os.FileInfo

func (a byChunk) Len() int      { return len(a) }
func (a byChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}
