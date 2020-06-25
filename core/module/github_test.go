// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	"os"
	"testing"

	"github.com/clivern/poodle/pkg"
)

// TestGithub test cases
func TestGithub(t *testing.T) {

	// export GITHUB_USERNAME=clivern
	// export GITHUB_OAUTH_TOKEN=~~secret~~
	if os.Getenv("GITHUB_USERNAME") == "" {
		return
	}

	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" {
		return
	}

	gistID := ""

	t.Run("TestGithubAuth", func(t *testing.T) {
		httpClient := NewHTTPClient()
		githubClient := NewGithubClient(
			httpClient,
			GithubAPI,
			os.Getenv("GITHUB_USERNAME"),
			os.Getenv("GITHUB_OAUTH_TOKEN"),
		)

		oauth, err := githubClient.Check(context.TODO())
		pkg.Expect(t, oauth.Scopes, "gist")
		pkg.Expect(t, oauth.Valid, true)
		pkg.Expect(t, err, nil)
	})

	t.Run("TestGithubGistCreate", func(t *testing.T) {
		httpClient := NewHTTPClient()
		githubClient := NewGithubClient(
			httpClient,
			GithubAPI,
			os.Getenv("GITHUB_USERNAME"),
			os.Getenv("GITHUB_OAUTH_TOKEN"),
		)

		result, err := githubClient.CreateGist(context.TODO(), Gist{
			Description: "Test Gist",
			Public:      false,
			Files: map[string]File{
				"config.toml": File{
					Content:  "some configs here",
					Filename: "config.toml",
				},
			},
		})

		gistID = result.ID
		pkg.Expect(t, result.ID != "", true)
		pkg.Expect(t, err, nil)
	})

	t.Run("TestGithubGistUpdate", func(t *testing.T) {
		httpClient := NewHTTPClient()
		githubClient := NewGithubClient(
			httpClient,
			GithubAPI,
			os.Getenv("GITHUB_USERNAME"),
			os.Getenv("GITHUB_OAUTH_TOKEN"),
		)

		result, err := githubClient.UpdateGist(context.TODO(), gistID, Gist{
			Description: "Update Test Gist",
			Public:      false,
			Files: map[string]File{
				"config.toml": File{
					Content:  "update configs here",
					Filename: "config.toml",
				},
			},
		})

		pkg.Expect(t, result.ID, gistID)
		pkg.Expect(t, err, nil)
	})

	t.Run("TestGithubGistGet", func(t *testing.T) {
		httpClient := NewHTTPClient()
		githubClient := NewGithubClient(
			httpClient,
			GithubAPI,
			os.Getenv("GITHUB_USERNAME"),
			os.Getenv("GITHUB_OAUTH_TOKEN"),
		)

		result, err := githubClient.GetGist(context.TODO(), gistID)

		pkg.Expect(t, result.ID, gistID)
		pkg.Expect(t, err, nil)
	})

	t.Run("TestGithubGistDelete", func(t *testing.T) {
		httpClient := NewHTTPClient()
		githubClient := NewGithubClient(
			httpClient,
			GithubAPI,
			os.Getenv("GITHUB_USERNAME"),
			os.Getenv("GITHUB_OAUTH_TOKEN"),
		)

		result, err := githubClient.DeleteGist(context.TODO(), gistID)

		pkg.Expect(t, result, true)
		pkg.Expect(t, err, nil)
	})
}
