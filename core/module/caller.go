// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

// Caller struct
type Caller struct {
	HTTPClient *HTTPClient
	APIURL     string
}

// NewCaller creates an instance of a caller
func NewCaller(httpClient *HTTPClient, apiURL string) Caller {
	client := Caller{}

	client.HTTPClient = httpClient
	client.APIURL = apiURL

	return client
}
