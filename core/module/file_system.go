// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clivern/poodle/core/util"

	"github.com/araddon/dateparse"
)

// VFile struct
type VFile struct {
	ModTimestamp int64
	Content      string
	Name         string
}

// FileSystem struct
type FileSystem struct {
	Files map[string]VFile
}

// NewFileSystem creates a file system object
func NewFileSystem() *FileSystem {
	return &FileSystem{
		Files: make(map[string]VFile),
	}
}

// LoadFromLocal gets files from a local path
func (f *FileSystem) LoadFromLocal(basePath string, extension string) error {

	f.Files = make(map[string]VFile)
	basePath = util.RemoveTrailingSlash(basePath)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		modifParsed, err := dateparse.ParseLocal(info.ModTime().String())

		if err != nil {
			return err
		}

		if basePath != path && !info.IsDir() && strings.Contains(path, extension) {
			file := strings.Replace(path, util.EnsureTrailingSlash(basePath), "", -1)

			content, err := util.ReadFile(path)

			if err != nil {
				return err
			}

			f.Files[file] = VFile{
				ModTimestamp: modifParsed.Unix(),
				Content:      content,
				Name:         file,
			}
		}

		return nil
	})

	return err
}

// DumpLocally updates a local path
func (f *FileSystem) DumpLocally(basePath string) error {

	for key, file := range f.Files {
		handler, err := os.Create(fmt.Sprintf(
			"%s%s",
			util.EnsureTrailingSlash(basePath),
			key,
		))

		if err != nil {
			return err
		}

		defer handler.Close()

		handler.WriteString(file.Content)
	}

	return nil
}

// Sync update fs object from remote fs
func (f *FileSystem) Sync(remoteFs *FileSystem) error {

	for key, file := range f.Files {
		if v, found := remoteFs.Files[key]; found {

			// Rewrite the remote file
			if file.ModTimestamp > v.ModTimestamp {
				remoteFs.Files[key] = file
			}

			// Rewrite the local file
			if file.ModTimestamp < v.ModTimestamp {
				f.Files[key] = file
			}

		} else {
			// File locally and missing remotely
			remoteFs.Files[key] = file
		}
	}

	for key, file := range remoteFs.Files {
		if v, found := f.Files[key]; found {

			// Rewrite the local file
			if file.ModTimestamp > v.ModTimestamp {
				f.Files[key] = file
			}

			// Rewrite the remote file
			if file.ModTimestamp < v.ModTimestamp {
				remoteFs.Files[key] = file
			}

		} else {
			// File remotely & not locally
			f.Files[key] = file
		}
	}

	return nil
}

// LoadFromJSON update object from json
func (f *FileSystem) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &f)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ConvertToJSON convert object to json
func (f *FileSystem) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
