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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a service definition file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var ok bool

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Delete command got called.")

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

		files, err := util.ListFiles(util.EnsureTrailingSlash(conf.Services.Directory))

		if err != nil {
			fmt.Printf(
				"Error while listing services under %s: %s",
				util.EnsureTrailingSlash(conf.Services.Directory),
				err.Error(),
			)
			return
		}

		data := []string{}
		index := map[string]string{}
		service := &model.Service{}

		for _, v := range files {
			if strings.Contains(v.Name, ".toml") {

				service = model.NewService(v.Name)
				err = service.Decode(v.Path)

				if err != nil {
					fmt.Printf(
						"Error while decoding service %s: %s",
						v.Path,
						err.Error(),
					)
					return
				}

				data = append(
					data,
					service.Main.ID,
				)

				index[service.Main.ID] = v.Path
			}
		}

		result := ""
		finder := module.FuzzyFinder{}
		prompt := module.Prompt{}

		if finder.Available() {
			result, err = finder.Show(data)
		} else {
			result, err = prompt.Select(
				fmt.Sprintf("Select a Service"),
				data,
			)
		}

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		choice, err := prompt.Select(
			fmt.Sprintf("Are you sure"),
			[]string{"No", "Yes"},
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		if choice == "No" {
			fmt.Println(Red("Service file deletion skipped!"))
			return
		}

		spin := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
		spin.Color("green")
		spin.Start()

		err = util.DeleteFile(index[result])

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		// Delete remotely
		if conf.Gist.Username == "" || conf.Gist.AccessToken == "" {
			spin.Stop()
			fmt.Println(Green("Service file deleted successfully!"))
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
			spin.Stop()
			fmt.Println(Green("Service file deleted successfully!"))
			return
		}

		if !oauth.Valid {
			spin.Stop()
			fmt.Println(Green("Service file deleted successfully!"))
			return
		}

		remoteGist, err := githubClient.GetGist(context.TODO(), conf.Gist.GistID)

		if err != nil {
			fmt.Printf(
				"Error while fetching remote gist: %s",
				err.Error(),
			)
			return
		}

		remoteFs := module.NewFileSystem()

		if _, found := remoteGist.Files["poodle"]; found {

			ok, err = remoteFs.LoadFromJSON([]byte(remoteGist.Files["poodle"].Content))

			if !ok || err != nil {
				fmt.Printf(
					"Error while loading remote fs from json: %s",
					err.Error(),
				)
				return
			}

			delete(remoteFs.Files, strings.Replace(
				index[result],
				util.EnsureTrailingSlash(conf.Services.Directory),
				"",
				-1,
			))

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

		spin.Stop()

		fmt.Println(Green("Service file deleted successfully!"))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
