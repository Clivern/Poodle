// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GithubAPI github api url
const GithubAPI = "https://api.github.com"

// Github struct
type Github struct {
	OAuth      OAuth
	HTTPClient *HTTPClient
	APIURL     string
}

// OAuth struct
type OAuth struct {
	Token    string
	Username string
	Scopes   string
	Valid    bool
}

// File struct
type File struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

// Gist struct
type Gist struct {
	Description string          `json:"description"`
	Public      bool            `json:"public"`
	Files       map[string]File `json:"files"`
}

// GistResponse struct
type GistResponse struct {
	ID string `json:"id"`
}

// NewGithubClient creates an instance of github client
func NewGithubClient(httpClient *HTTPClient, apiURL, username, token string) Github {
	client := Github{}

	client.HTTPClient = httpClient
	client.APIURL = apiURL

	client.OAuth = OAuth{
		Token:    token,
		Username: username,
	}

	return client
}

// Check fetch current OAuth token data
func (g *Github) Check(ctx context.Context) (OAuth, error) {
	response, err := g.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/users/%s", g.APIURL, g.OAuth.Username),
		map[string]string{},
		map[string]string{"Authorization": fmt.Sprintf("token %s", g.OAuth.Token)},
	)

	if err != nil {
		return g.OAuth, err
	}

	if http.StatusOK != g.HTTPClient.GetStatusCode(response) {
		return g.OAuth, fmt.Errorf("Invalid status code %d", g.HTTPClient.GetStatusCode(response))
	}

	g.OAuth.Scopes = g.HTTPClient.GetHeaderValue(response, "X-OAuth-Scopes")
	g.OAuth.Valid = strings.Contains(g.OAuth.Scopes, "gist")

	if !g.OAuth.Valid {
		return g.OAuth, fmt.Errorf("Gist scope not allowed")
	}

	return g.OAuth, nil
}

// CreateGist creates a gist
func (g *Github) CreateGist(ctx context.Context, gist Gist) (GistResponse, error) {
	gistResponse := GistResponse{}

	request, err := gist.ConvertToJSON()

	if err != nil {
		return gistResponse, err
	}

	response, err := g.HTTPClient.Post(
		ctx,
		fmt.Sprintf("%s/gists", g.APIURL),
		request,
		map[string]string{},
		map[string]string{"Authorization": fmt.Sprintf("token %s", g.OAuth.Token)},
	)

	if err != nil {
		return gistResponse, err
	}

	if http.StatusCreated != g.HTTPClient.GetStatusCode(response) {
		return gistResponse, fmt.Errorf("Invalid status code %d", g.HTTPClient.GetStatusCode(response))
	}

	body, err := g.HTTPClient.ToString(response)

	if err != nil {
		return gistResponse, err
	}

	gistResponse.LoadFromJSON([]byte(body))

	return gistResponse, nil
}

// UpdateGist updates a gist
func (g *Github) UpdateGist(ctx context.Context, id string, gist Gist) (GistResponse, error) {
	gistResponse := GistResponse{}

	request, err := gist.ConvertToJSON()

	if err != nil {
		return gistResponse, err
	}

	response, err := g.HTTPClient.Patch(
		ctx,
		fmt.Sprintf("%s/gists/%s", g.APIURL, id),
		request,
		map[string]string{},
		map[string]string{"Authorization": fmt.Sprintf("token %s", g.OAuth.Token)},
	)

	if err != nil {
		return gistResponse, err
	}

	if http.StatusOK != g.HTTPClient.GetStatusCode(response) {
		return gistResponse, fmt.Errorf("Invalid status code %d", g.HTTPClient.GetStatusCode(response))
	}

	body, err := g.HTTPClient.ToString(response)

	if err != nil {
		return gistResponse, err
	}

	gistResponse.LoadFromJSON([]byte(body))

	return gistResponse, nil
}

// GetGist gets a gist
func (g *Github) GetGist(ctx context.Context, id string) (GistResponse, error) {
	gistResponse := GistResponse{}

	response, err := g.HTTPClient.Get(
		ctx,
		fmt.Sprintf("%s/gists/%s", g.APIURL, id),
		map[string]string{},
		map[string]string{"Authorization": fmt.Sprintf("token %s", g.OAuth.Token)},
	)

	if err != nil {
		return gistResponse, err
	}

	if http.StatusOK != g.HTTPClient.GetStatusCode(response) {
		return gistResponse, fmt.Errorf("Invalid status code %d", g.HTTPClient.GetStatusCode(response))
	}

	body, err := g.HTTPClient.ToString(response)

	if err != nil {
		return gistResponse, err
	}

	gistResponse.LoadFromJSON([]byte(body))

	return gistResponse, nil
}

// DeleteGist deletes a gist
func (g *Github) DeleteGist(ctx context.Context, id string) (bool, error) {
	response, err := g.HTTPClient.Delete(
		ctx,
		fmt.Sprintf("%s/gists/%s", g.APIURL, id),
		map[string]string{},
		map[string]string{"Authorization": fmt.Sprintf("token %s", g.OAuth.Token)},
	)

	if err != nil {
		return false, err
	}

	if http.StatusNoContent != g.HTTPClient.GetStatusCode(response) {
		return false, fmt.Errorf("Invalid status code %d", g.HTTPClient.GetStatusCode(response))
	}

	return true, nil
}

// LoadFromJSON update object from json
func (g *Gist) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &g)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (g *Gist) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&g)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (g *GistResponse) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &g)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (g *GistResponse) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&g)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
