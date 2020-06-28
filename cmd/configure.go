// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/clivern/poodle/core/model"
	"github.com/clivern/poodle/core/module"
	"github.com/clivern/poodle/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Poodle",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Configure command got called.")

		file := fmt.Sprintf(
			"%s%s",
			util.EnsureTrailingSlash(os.Getenv("HOME")),
			"poodle/config.toml",
		)

		log.WithFields(log.Fields{
			"file": file,
		}).Debug("Create config file if not exists.")

		if !util.FileExists(file) {
			log.WithFields(log.Fields{
				"file": file,
			}).Debug("Creating config file")

			err = util.StoreFile(file, "")
		}

		if err != nil {
			fmt.Printf(
				"Error while creating file %s: %s",
				file,
				err.Error(),
			)
			return
		}

		conf := model.NewConfigs()
		err = conf.Decode(file)

		if err != nil {
			fmt.Printf(
				"Error while decoding configs %s: %s",
				file,
				err.Error(),
			)
			return
		}

		err = conf.Encode(file)

		if err != nil {
			fmt.Printf(
				"Error while encoding configs %s: %s",
				file,
				err.Error(),
			)
			return
		}

		prompt := module.Prompt{}

		log.Debug("Start interactive prompt")

		confWith, err := prompt.Select(
			fmt.Sprintf("Configure With"),
			[]string{"Interactive", fmt.Sprintf("Editor (%s)", os.Getenv("EDITOR"))},
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		log.WithFields(log.Fields{
			"confWith": confWith,
		}).Debug("Interactive prompt")

		if confWith != "Interactive" {
			editor := module.Editor{}
			err = editor.Edit(file)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				log.WithFields(log.Fields{
					"file": file,
				}).Debug("Configs Updated")

				fmt.Println("Configs Updated")
			}
			return
		}

		username, err := prompt.Input(
			fmt.Sprintf("Github Username:"),
			module.NotEmpty,
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		// Override github username
		conf.Gist.Username = username

		token, err := prompt.Input(
			fmt.Sprintf("Github OAuth Token:"),
			module.NotEmpty,
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		// Override github token
		conf.Gist.AccessToken = token

		servicesDir, err := prompt.Input(
			fmt.Sprintf("Services Definitions Directory:"),
			module.NotEmpty,
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		// Override services dir
		conf.Services.Directory = servicesDir

		err = conf.Encode(file)

		if err != nil {
			fmt.Printf(
				"Error while encoding configs %s: %s",
				file,
				err.Error(),
			)
			return
		}

		log.WithFields(log.Fields{
			"file": file,
		}).Debug("Configs Updated")

		fmt.Println("Configs Updated")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
