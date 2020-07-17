// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"
	"regexp"
	"strings"

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

	for _, end := range service.Endpoint {
		if fmt.Sprintf("%s - %s", service.Main.ID, end.ID) != endpointID {
			continue
		}
		fields = c.MergeFields(fields, c.ParseFields(end.URI))
		fields = c.MergeFields(fields, c.ParseFields(end.Body))
	}

	return fields
}

// ParseFields parses a string to fetch fields
func (c *Caller) ParseFields(data string) map[string]Field {
	var ita []string
	m := regexp.MustCompile(`{\$(.*?)}`)
	items := m.FindAllString(data, -1)
	fields := make(map[string]Field)

	for _, item := range items {
		item = strings.Replace(item, "$", "", -1)
		item = strings.Replace(item, "{", "", -1)
		item = strings.Replace(item, "}", "", -1)

		if strings.Contains(item, ":") {
			ita = strings.Split(item, ":")
			fields[ita[0]] = Field{
				Prompt:     fmt.Sprintf(`$%s (default='%s'):`, ita[0], ita[1]),
				IsOptional: true,
				Default:    ita[1],
			}
		} else {
			fields[item] = Field{
				Prompt:     fmt.Sprintf(`$%s (default=''):`, item),
				IsOptional: false,
				Default:    "",
			}
		}
	}

	return fields
}

// MergeFields merges two fields list
func (c *Caller) MergeFields(m1, m2 map[string]Field) map[string]Field {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

// Call calls the remote service
func (c *Caller) Call(endpointID string, service *model.Service, fields map[string]Field) string {
	fmt.Println(endpointID)
	fmt.Println(fields)

	return ""
}
