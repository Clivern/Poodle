// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/clivern/poodle/core/model"
	"github.com/clivern/poodle/core/module"
	"github.com/clivern/poodle/core/util"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync API services definitions",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var ok bool

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		spin := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		spin.Color("green")
		spin.Start()

		log.Debug("Sync command got called.")

		if !util.FileExists(Config) {
			fmt.Printf(
				"Config file is missing %s, Please start with $ poodle configure",
				Config,
			)
			return
		}

		conf := model.NewConfigs()
		err = conf.Decode(Config)

		if err != nil {
			fmt.Printf(
				"Error while decoding configs %s: %s",
				Config,
				err.Error(),
			)
			return
		}

		githubClient := module.NewGithubClient(
			module.NewHTTPClient(),
			module.GithubAPI,
			conf.Gist.Username,
			conf.Gist.AccessToken,
		)

		oauth, err := githubClient.Check(context.TODO())

		if err != nil {
			fmt.Printf(
				"Error during github authentication for username %s: %s",
				conf.Gist.Username,
				err.Error(),
			)
			return
		}

		if !oauth.Valid {
			fmt.Printf(
				"Invalid github username, token or scopes [%s] don't include gist",
				oauth.Scopes,
			)
			return
		}

		found := true

		// Validate if gist still exist
		if strings.TrimSpace(conf.Gist.GistID) != "" {
			res, err := githubClient.GetGist(context.TODO(), conf.Gist.GistID)

			if err != nil || res.ID == "" {
				found = false
			}
		}

		if strings.TrimSpace(conf.Gist.GistID) == "" || !found {
			// Create github gist
			files := make(map[string]module.File)
			files["poodle"] = module.File{
				Content:  "#",
				Filename: "poodle",
			}

			result, err := githubClient.CreateGist(context.TODO(), module.Gist{
				Description: "Poodle",
				Public:      conf.Gist.Public,
				Files:       files,
			})

			if err != nil {
				fmt.Printf(
					"Error while creating a github gist: %s",
					err.Error(),
				)
				return
			}

			conf.Gist.GistID = result.ID

			err = conf.Encode(Config)

			if err != nil {
				fmt.Printf(
					"Error while encoding configs %s: %s",
					Config,
					err.Error(),
				)
				return
			}
		}

		status, err := githubClient.GetSyncStatus(
			context.TODO(),
			conf.Services.Directory,
			conf.Gist.GistID,
		)

		if err != nil {
			fmt.Printf(
				"Error while fetching sync status: %s",
				err.Error(),
			)
			return
		}

		if status == "upload" {
			ok, err = githubClient.SyncByUpload(
				context.TODO(),
				conf.Services.Directory,
				conf.Gist.GistID,
			)
		} else {
			ok, err = githubClient.SyncByDownload(
				context.TODO(),
				conf.Services.Directory,
				conf.Gist.GistID,
			)
		}

		if !ok || err != nil {
			fmt.Printf(
				"Error while syncing services definitions directory: %s",
				err.Error(),
			)
			return
		}

		log.WithFields(log.Fields{
			"status": status,
		}).Debug("Sync Done")

		spin.Stop()

		fmt.Println("Already up-to-date")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
