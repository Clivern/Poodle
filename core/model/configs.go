// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"os"

	"github.com/clivern/poodle/core/util"

	"github.com/BurntSushi/toml"
)

// Configs type
type Configs struct {
	General  General  `toml:"General"`
	Gist     Gist     `toml:"Gist"`
	Services Services `toml:"Services"`
}

// General type
type General struct {
	Editor    string `toml:"editor"`
	Column    int    `toml:"column"`
	Selectcmd string `toml:"selectcmd"`
	Backend   string `toml:"backend"`
	Sortby    string `toml:"sortby"`
}

// Gist type
type Gist struct {
	AccessToken string `toml:"access_token"`
	GistID      string `toml:"gist_id"`
	Username    string `toml:"username"`
	Public      bool   `toml:"public"`
	AutoSync    bool   `toml:"auto_sync"`
}

// Services type
type Services struct {
	Directory string `toml:"directory"`
}

// NewConfigs creates an instance of Configs
func NewConfigs() *Configs {
	path := fmt.Sprintf(
		"%s%s",
		util.EnsureTrailingSlash(os.Getenv("HOME")),
		"poodle/definitions",
	)

	return &Configs{
		General: General{
			Editor:    "nano",
			Column:    40,
			Selectcmd: "fzf --ansi",
			Backend:   "gist",
		},
		Services: Services{
			Directory: path,
		},
	}
}

// Decode decodes from file to struct
func (g *Configs) Decode(path string) error {
	if _, err := toml.DecodeFile(path, &g); err != nil {
		return err
	}

	return nil
}

// Encode encodes struct and store on file
func (g *Configs) Encode(path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	err = toml.NewEncoder(f).Encode(g)

	if err != nil {
		return err
	}

	return nil
}
