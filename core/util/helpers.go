// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/satori/go.uuid"
)

// File struct
type File struct {
	Path         string
	Name         string
	Size         int64
	ModTime      time.Time
	ModTimestamp int64
}

// InArray check if value is on array
func InArray(val interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

// GenerateUUID4 create a UUID
func GenerateUUID4() string {
	u := uuid.Must(uuid.NewV4(), nil)
	return u.String()
}

// ListFiles lists all files inside a dir
func ListFiles(basePath string) ([]File, error) {
	var files []File

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		modifParsed, err := dateparse.ParseLocal(info.ModTime().String())

		if err != nil {
			return err
		}

		if basePath != path && !info.IsDir() {
			files = append(files, File{
				Path:         path,
				Name:         info.Name(),
				Size:         info.Size(),
				ModTime:      info.ModTime(),
				ModTimestamp: modifParsed.Unix(),
			})
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}

// ReadFile get the file content
func ReadFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FilterFiles filters files list based on specific sub-strings
func FilterFiles(files, filters []string) []string {
	var filteredFiles []string

	for _, file := range files {
		ok := true
		for _, filter := range filters {

			ok = ok && strings.Contains(file, filter)
		}
		if ok {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

// Unset remove element at position i
func Unset(a []string, i int) []string {
	a[i] = a[len(a)-1]
	a[len(a)-1] = ""
	return a[:len(a)-1]
}

// ConvertToJSON convert object to json
func ConvertToJSON(val interface{}) (string, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
