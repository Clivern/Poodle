// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Basic type
type Basic struct {
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	Header   []string `toml:"header"`
}

// APIKey type
type APIKey struct {
	Header []string `toml:"header"`
}

// Bearer type
type Bearer struct {
	Header []string `toml:"header"`
}

// Security type
type Security struct {
	Scheme string `toml:"scheme"`
	Basic  Basic  `toml:"Basic"`
	APIKey APIKey `toml:"ApiKey"`
	Bearer Bearer `toml:"Bearer"`
}

// Main type
type Main struct {
	ID          string     `toml:"id"`
	Name        string     `toml:"name"`
	Description string     `toml:"description"`
	Timeout     string     `toml:"timeout"`
	ServiceURL  string     `toml:"service_url"`
	Headers     [][]string `toml:"headers"`
}

// Endpoint type
type Endpoint struct {
	ID          string     `toml:"id"`
	Name        string     `toml:"name"`
	Description string     `toml:"description"`
	Method      string     `toml:"method"`
	Headers     [][]string `toml:"headers"`
	Parameters  [][]string `toml:"parameters"`
	URI         string     `toml:"uri"`
	Body        string     `toml:"body"`
}

// Service type
type Service struct {
	Main     Main       `toml:"Main"`
	Security Security   `toml:"Security"`
	Endpoint []Endpoint `toml:"Endpoint"`
}

// NewService creates an instance of Service
func NewService(id string) *Service {
	return &Service{
		Main: Main{
			ID:          id,
			Timeout:     "30s",
			Name:        "~~ name Here ~~",
			Description: "~~ description Here ~~",
			ServiceURL:  "https://httpbin.org",
			Headers: [][]string{
				[]string{"Content-Type", "application/json"},
			},
		},
		Security: Security{
			Scheme: "none",
			Basic: Basic{
				Username: "{$authUsername:default}",
				Password: "{$authPassword:default}",
				Header: []string{
					"Authorization",
					"Basic base64(username:password)",
				},
			},
			APIKey: APIKey{
				Header: []string{
					"X-API-KEY",
					"{$authApiKey:default}",
				},
			},
			Bearer: Bearer{
				Header: []string{
					"Authorization",
					"Bearer {$authBearerToken:default}",
				},
			},
		},
		Endpoint: []Endpoint{
			Endpoint{
				ID:          "CreateItem",
				Name:        "Create an item",
				Description: "",
				Method:      "post",
				Headers:     [][]string{},
				Parameters:  [][]string{},
				URI:         "/anything",
				Body:        `{"name":"{$name}","type":"{$type:default}"}`,
			},
			Endpoint{
				ID:          "GetItems",
				Name:        "Get a list of items",
				Description: "",
				Method:      "get",
				Headers:     [][]string{},
				Parameters: [][]string{
					[]string{"limit", "{$limit:100}"},
					[]string{"offset", "{$offset:0}"},
				},
				URI:  "/anything",
				Body: "",
			},
			Endpoint{
				ID:          "GetItem",
				Name:        "Get an items",
				Description: "",
				Method:      "get",
				Headers:     [][]string{},
				Parameters:  [][]string{},
				URI:         "/anything/{$id}",
				Body:        "",
			},
			Endpoint{
				ID:          "UpdateItem",
				Name:        "Update an item",
				Description: "",
				Method:      "put",
				Headers:     [][]string{},
				Parameters:  [][]string{},
				URI:         "/anything/{$id}",
				Body:        `{"name":"{$name}","type":"{$type:default}"}`,
			},
			Endpoint{
				ID:          "DeleteItem",
				Name:        "Delete an item",
				Description: "",
				Method:      "delete",
				Headers:     [][]string{},
				Parameters:  [][]string{},
				URI:         "/anything/{$id}",
				Body:        "",
			},
		},
	}
}

// NewEmptyService creates an empty instance of Service
func NewEmptyService(id string) *Service {
	return &Service{
		Main: Main{
			ID:          id,
			Timeout:     "30s",
			Name:        "",
			Description: "",
			ServiceURL:  "",
			Headers:     [][]string{},
		},
		Security: Security{
			Scheme: "none",
			Basic: Basic{
				Username: "{$authUsername:default}",
				Password: "{$authPassword:default}",
				Header: []string{
					"Authorization",
					"Basic base64(username:password)",
				},
			},
			APIKey: APIKey{
				Header: []string{
					"X-API-KEY",
					"{$authApiKey:default}",
				},
			},
			Bearer: Bearer{
				Header: []string{
					"Authorization",
					"Bearer {$authBearerToken:default}",
				},
			},
		},
		Endpoint: []Endpoint{},
	}
}

// Decode decodes from file to struct
func (s *Service) Decode(path string) error {
	if _, err := toml.DecodeFile(path, &s); err != nil {
		return err
	}

	return nil
}

// Encode encodes struct and store on file
func (s *Service) Encode(path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	err = toml.NewEncoder(f).Encode(s)

	if err != nil {
		return err
	}

	return nil
}
