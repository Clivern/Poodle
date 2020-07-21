// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// File struct
type File struct {
	Path         string
	Ident        string
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

// ListFiles lists all files inside a dir
func ListFiles(basePath string) (map[string]File, error) {
	ident := ""
	files := make(map[string]File)
	basePath = RemoveTrailingSlash(basePath)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		modifParsed, err := dateparse.ParseLocal(info.ModTime().String())

		if err != nil {
			return err
		}

		if basePath != path && !info.IsDir() {
			ident = strings.Replace(path, EnsureTrailingSlash(basePath), "", -1)
			ident = strings.Replace(ident, string(os.PathSeparator), "__", -1)
			files[ident] = File{
				Path:         path,
				Ident:        strings.Replace(path, string(os.PathSeparator), "__", -1),
				Name:         info.Name(),
				Size:         info.Size(),
				ModTime:      info.ModTime(),
				ModTimestamp: modifParsed.Unix(),
			}
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

// EnsureTrailingSlash ensure there is a trailing slash
func EnsureTrailingSlash(dir string) string {
	return fmt.Sprintf(
		"%s%s",
		strings.TrimRight(dir, string(os.PathSeparator)),
		string(os.PathSeparator),
	)
}

// RemoveTrailingSlash removes any trailing slash
func RemoveTrailingSlash(dir string) string {
	return strings.TrimRight(dir, string(os.PathSeparator))
}

// RemoveStartingSlash removes any starting slash
func RemoveStartingSlash(dir string) string {
	return strings.TrimLeft(dir, string(os.PathSeparator))
}

// ClearDir removes all files and sub dirs
func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// StoreFile stores a file content
func StoreFile(path, content string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0775)

	if err != nil {
		return err
	}

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)

	return err
}

// PathExists reports whether the path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileExists reports whether the named file exists
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// DirExists reports whether the dir exists
func DirExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

// EnsureDir ensures that directory exists
func EnsureDir(dirName string, mode int) (bool, error) {
	err := os.MkdirAll(dirName, os.FileMode(mode))

	if err == nil || os.IsExist(err) {
		return true, nil
	}
	return false, err
}
