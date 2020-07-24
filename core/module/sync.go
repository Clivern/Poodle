// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

// Item struct
type Item struct {
	Timestamp int
	Content string
}

// Sync struct
type Sync struct {
}

// CollectFs get
func (s *Sync) CollectFs(path string) (string, error) {

}

// LoadFromJSON update object from json
func (s *Sync) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ConvertToJSON convert object to json
func (s *Sync) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
