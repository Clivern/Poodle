// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/poodle/core/model"
	"github.com/clivern/poodle/core/util"
	. "github.com/logrusorgru/aurora/v3"
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

		// Get api key if auth is api_key
		if service.Security.Scheme == "api_key" {
			fields = c.MergeFields(fields, c.ParseFields(service.Security.APIKey.Header[1]))
		}

		// Get bearer token if auth is bearer
		if service.Security.Scheme == "bearer" {
			fields = c.MergeFields(fields, c.ParseFields(service.Security.Bearer.Header[1]))
		}

		// Get username and password if auth is basic
		if service.Security.Scheme == "basic" {
			fields = c.MergeFields(fields, c.ParseFields(service.Security.Basic.Username))
			fields = c.MergeFields(fields, c.ParseFields(service.Security.Basic.Password))
		}

		// Get URI vars
		fields = c.MergeFields(fields, c.ParseFields(end.URI))

		// Get headers vars
		for _, header := range end.Headers {
			fields = c.MergeFields(fields, c.ParseFields(header[1]))
		}

		// Get parameters vars
		for _, parameter := range end.Parameters {
			fields = c.MergeFields(fields, c.ParseFields(parameter[1]))
		}

		// Get Body vars
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
				Prompt:     fmt.Sprintf(`$%s (default='%s'):`, ita[0], Yellow(ita[1])),
				IsOptional: true,
				Default:    ita[1],
			}
		} else {
			fields[item] = Field{
				Prompt:     fmt.Sprintf(`$%s%s (default=''):`, item, Red("*")),
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
func (c *Caller) Call(endpointID string, service *model.Service, fields map[string]Field) (*http.Response, error) {
	var response *http.Response
	var err error

	for _, end := range service.Endpoint {
		if fmt.Sprintf("%s - %s", service.Main.ID, end.ID) != endpointID {
			continue
		}

		url := fmt.Sprintf(
			"%s%s",
			util.EnsureTrailingSlash(service.Main.ServiceURL),
			util.RemoveStartingSlash(c.ReplaceVars(end.URI, fields)),
		)

		data := c.ReplaceVars(end.Body, fields)
		parameters := make(map[string]string)
		headers := make(map[string]string)

		// addd service global headers
		for _, header := range service.Main.Headers {
			headers[header[0]] = header[1]
		}

		// Add api key to headers if auth is api_key
		if service.Security.Scheme == "api_key" {
			headers[service.Security.APIKey.Header[0]] = c.ReplaceVars(service.Security.APIKey.Header[1], fields)
		}

		// Add bearer token to headers if auth is bearer
		if service.Security.Scheme == "bearer" {
			headers[service.Security.Bearer.Header[0]] = c.ReplaceVars(service.Security.Bearer.Header[1], fields)
		}

		// Add base64 of username & password if auth is basic
		if service.Security.Scheme == "basic" {
			username := c.ReplaceVars(service.Security.Basic.Username, fields)
			password := c.ReplaceVars(service.Security.Basic.Password, fields)

			headers[service.Security.Basic.Header[0]] = strings.Replace(
				service.Security.Basic.Header[1],
				"base64(username:password)",
				b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
				-1,
			)
		}

		// Get headers vars
		for _, header := range end.Headers {
			headers[header[0]] = c.ReplaceVars(header[1], fields)
		}

		// Get parameters vars
		for _, parameter := range end.Parameters {
			parameters[parameter[0]] = c.ReplaceVars(parameter[1], fields)
		}

		timeout, err := strconv.Atoi(strings.Replace(service.Main.Timeout, "s", "", -1))

		if err != nil {
			return response, err
		}

		c.HTTPClient.Timeout = time.Duration(timeout)

		if end.Method == "get" {
			response, err = c.HTTPClient.Get(
				context.TODO(),
				url,
				parameters,
				headers,
			)
		}

		if end.Method == "post" {
			response, err = c.HTTPClient.Post(
				context.TODO(),
				url,
				data,
				parameters,
				headers,
			)
		}

		if end.Method == "put" {
			response, err = c.HTTPClient.Put(
				context.TODO(),
				url,
				data,
				parameters,
				headers,
			)
		}

		if end.Method == "delete" {
			response, err = c.HTTPClient.Delete(
				context.TODO(),
				url,
				parameters,
				headers,
			)
		}

		if end.Method == "patch" {
			response, err = c.HTTPClient.Patch(
				context.TODO(),
				url,
				data,
				parameters,
				headers,
			)
		}
	}

	if err != nil {
		return response, err
	}

	return response, nil
}

// ReplaceVars replaces vars
func (c *Caller) ReplaceVars(data string, fields map[string]Field) string {
	for k, field := range fields {
		if field.IsOptional {
			data = strings.Replace(
				data,
				fmt.Sprintf("{$%s:%s}", k, field.Default),
				field.Value,
				-1,
			)
		} else {
			data = strings.Replace(
				data,
				fmt.Sprintf("{$%s}", k),
				field.Value,
				-1,
			)
		}
	}

	return data
}

// Pretty returns colored output
func (c *Caller) Pretty(response *http.Response) string {
	body, err := c.HTTPClient.ToString(response)

	if err != nil {
		body = fmt.Sprintf("Error %s", err.Error())
	}

	responseCode := c.HTTPClient.GetStatusCode(response)

	value := "\n---\n"

	value = value + fmt.Sprintf(
		"%s %d %s\n",
		Blue(response.Proto),
		Blue(responseCode),
		Cyan(http.StatusText(responseCode)),
	)

	for k, v := range response.Header {
		for _, h := range v {
			value = value + fmt.Sprintf("%s: %s\n", Cyan(k), h)
		}
	}

	value = value + fmt.Sprintf("\n%s", Yellow(body))

	return value
}
