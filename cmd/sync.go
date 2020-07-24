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
	. "github.com/logrusorgru/aurora/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync services definitions",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var ok bool

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		spin := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
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

		localFS := module.NewFileSystem()
		remoteFs := module.NewFileSystem()

		err = localFS.LoadFromLocal(conf.Services.Directory, "toml")

		if err != nil {
			fmt.Printf(
				"Error while reading local files inside %s: %s",
				conf.Services.Directory,
				err.Error(),
			)
			return
		}

		localData, err := localFS.ConvertToJSON()

		if err != nil {
			fmt.Printf(
				"Error while converting local files data to json %s: %s",
				conf.Services.Directory,
				err.Error(),
			)
			return
		}

		if strings.TrimSpace(conf.Gist.GistID) == "" || !found {
			// Create github gist
			files := make(map[string]module.File)
			files["poodle"] = module.File{
				Content:  localData,
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

		remoteGist, err := githubClient.GetGist(context.TODO(), conf.Gist.GistID)

		if err != nil {
			fmt.Printf(
				"Error while fetching remote gist: %s",
				err.Error(),
			)
			return
		}

		if _, found := remoteGist.Files["poodle"]; found {

			ok, err = remoteFs.LoadFromJSON([]byte(remoteGist.Files["poodle"].Content))

			if !ok || err != nil {
				fmt.Printf(
					"Error while loading remote fs from json: %s",
					err.Error(),
				)
				return
			}

			err = localFS.Sync(remoteFs)

			if err != nil {
				fmt.Printf(
					"Error while sync remote and local fs: %s",
					err.Error(),
				)
				return
			}

			remoteData, err := remoteFs.ConvertToJSON()

			if err != nil {
				fmt.Printf(
					"Error while converting remote data to json: %s",
					err.Error(),
				)
				return
			}

			remoteGist.Files["poodle"] = module.File{
				Content:  remoteData,
				Filename: "poodle",
			}
		} else {
			remoteGist.Files["poodle"] = module.File{
				Content:  localData,
				Filename: "poodle",
			}
		}

		fmt.Println(localFS.Files)

		err = localFS.DumpLocally(conf.Services.Directory)

		if err != nil {
			fmt.Printf(
				"Error while updating local files: %s",
				err.Error(),
			)
			return
		}

		_, err = githubClient.UpdateGist(
			context.TODO(),
			conf.Gist.GistID,
			module.Gist{
				Description: "Poodle",
				Public:      conf.Gist.Public,
				Files:       remoteGist.Files,
			},
		)

		if err != nil {
			fmt.Printf(
				"Error while updating remote gist: %s",
				err.Error(),
			)
			return
		}

		log.Debug("Sync Done")

		spin.Stop()

		fmt.Println(Green("Already up-to-date"))
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
