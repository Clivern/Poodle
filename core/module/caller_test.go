// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"
	"strings"
	"testing"

	"github.com/clivern/poodle/core/model"

	"github.com/clivern/poodle/pkg"
)

// TestCallerPostRequest test cases
func TestCallerPostRequest(t *testing.T) {
	t.Run("TestCallerPostRequest", func(t *testing.T) {
		httpClient := NewHTTPClient()
		caller := NewCaller(httpClient)
		service := model.NewService("anything")

		fields01 := caller.GetFields(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[1].ID),
			service,
		)

		pkg.Expect(t, true, strings.Contains(fields01["name"].Prompt, "default"))
		pkg.Expect(t, false, fields01["name"].IsOptional)
		pkg.Expect(t, "", fields01["name"].Default)
		pkg.Expect(t, "", fields01["name"].Value)

		pkg.Expect(t, true, strings.Contains(fields01["type"].Prompt, "default"))
		pkg.Expect(t, true, fields01["type"].IsOptional)
		pkg.Expect(t, "default", fields01["type"].Default)
		pkg.Expect(t, "", fields01["type"].Value)

		res, err := caller.Call(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[1].ID),
			service,
			fields01,
		)

		pkg.Expect(t, nil, err)

		body, err := httpClient.ToString(res)

		pkg.Expect(t, nil, err)
		pkg.Expect(t, 200, httpClient.GetStatusCode(res))
		pkg.Expect(t, true, strings.Contains(body, "POST"))
	})
}

// TestCallerGetRequest test cases
func TestCallerGetRequest(t *testing.T) {
	t.Run("TestCallerGetRequest", func(t *testing.T) {
		httpClient := NewHTTPClient()
		caller := NewCaller(httpClient)
		service := model.NewService("anything")

		service.Endpoint[3].URI = "/anything"

		fields01 := caller.GetFields(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[3].ID),
			service,
		)

		res, err := caller.Call(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[3].ID),
			service,
			fields01,
		)

		pkg.Expect(t, nil, err)

		body, err := httpClient.ToString(res)

		pkg.Expect(t, nil, err)
		pkg.Expect(t, 200, httpClient.GetStatusCode(res))
		pkg.Expect(t, true, strings.Contains(body, "GET"))
	})
}

// TestCallerPutRequest test cases
func TestCallerPutRequest(t *testing.T) {
	t.Run("TestCallerPutRequest", func(t *testing.T) {
		httpClient := NewHTTPClient()
		caller := NewCaller(httpClient)
		service := model.NewService("anything")

		service.Endpoint[4].URI = "/anything"

		fields01 := caller.GetFields(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[4].ID),
			service,
		)

		res, err := caller.Call(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[4].ID),
			service,
			fields01,
		)

		pkg.Expect(t, nil, err)

		body, err := httpClient.ToString(res)

		pkg.Expect(t, nil, err)
		pkg.Expect(t, 200, httpClient.GetStatusCode(res))
		pkg.Expect(t, true, strings.Contains(body, "PUT"))
	})
}

// TestCallerDeleteRequest test cases
func TestCallerDeleteRequest(t *testing.T) {
	t.Run("TestCallerDeleteRequest", func(t *testing.T) {
		httpClient := NewHTTPClient()
		caller := NewCaller(httpClient)
		service := model.NewService("anything")

		service.Endpoint[5].URI = "/anything"

		fields01 := caller.GetFields(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[5].ID),
			service,
		)

		res, err := caller.Call(
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[5].ID),
			service,
			fields01,
		)

		pkg.Expect(t, nil, err)

		body, err := httpClient.ToString(res)

		pkg.Expect(t, nil, err)
		pkg.Expect(t, 200, httpClient.GetStatusCode(res))
		pkg.Expect(t, true, strings.Contains(body, "DELETE"))
	})
}
