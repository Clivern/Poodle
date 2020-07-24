// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

// File struct
type File struct {
	Timestamp int
	Content   string
}

// FileSystem struct
type FileSystem struct {
}

// LoadFromLocal gets files from a local path
func (f *FileSystem) LoadFromLocal(basePath string) error {

}

// DumpLocally updates a local path
func (f *FileSystem) DumpLocally(basePath string) error {

}

// UpdateFromRemote update fs object from remote fs
func (f *FileSystem) UpdateFromRemote(remoteFs FileSystem) error {

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
