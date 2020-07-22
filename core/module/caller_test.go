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
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[0].ID),
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
			fmt.Sprintf("%s - %s", service.Main.ID, service.Endpoint[0].ID),
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
