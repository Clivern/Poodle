// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"

	"github.com/clivern/poodle/core/model"
)

// Caller struct
type Caller struct {
	HTTPClient *HTTPClient
}

// Field struct
type Field struct {
	Prompt     string
	IsOptional bool
	Default    string
	Value      string
}

// NewCaller creates an instance of a caller
func NewCaller(httpClient *HTTPClient) Caller {
	client := Caller{}
	client.HTTPClient = httpClient

	return client
}

// GetFields get fields to collect from end user
func (c *Caller) GetFields(endpointID string, service *model.Service) map[string]Field {
	fields := make(map[string]Field)

	fields["id"] = Field{
		Prompt:     `$id (default=''):`,
		IsOptional: true,
		Default:    "",
	}

	fields["key"] = Field{
		Prompt:     `$key (default='def_key'):`,
		IsOptional: true,
		Default:    "def_key",
	}

	fields["name"] = Field{
		Prompt:     `$name (default=''):`,
		IsOptional: false,
		Default:    "",
	}

	return fields
}

// Call calls the remote service
func (c *Caller) Call(endpointID string, service *model.Service, fields map[string]Field) string {
	fmt.Println(endpointID)
	fmt.Println(fields)

	return ""
}
