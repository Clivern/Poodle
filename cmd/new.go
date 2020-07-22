// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/clivern/poodle/core/model"
	"github.com/clivern/poodle/core/module"
	"github.com/clivern/poodle/core/util"

	. "github.com/logrusorgru/aurora/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new service definition file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("New command got called.")

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

		prompt := module.Prompt{}

		relPath, err := prompt.Input(
			fmt.Sprintf("Service Id:"),
			func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("Input must not be empty")
				}

				match, err := regexp.MatchString("^[A-Za-z0-9-_/]+$", input)

				if !match || err != nil {
					return fmt.Errorf("Service Id must be alphanumeric")
				}

				path := fmt.Sprintf(
					"%s%s.toml",
					util.EnsureTrailingSlash(conf.Services.Directory),
					input,
				)

				if util.FileExists(path) {
					return fmt.Errorf("Service Id is used before")
				}

				return nil
			},
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		absPath := fmt.Sprintf(
			"%s%s.toml",
			util.EnsureTrailingSlash(conf.Services.Directory),
			relPath,
		)

		service := model.NewService(relPath)
		err = service.Encode(absPath)

		if err != nil {
			fmt.Printf(
				"Error while decoding configs %s: %s",
				absPath,
				err.Error(),
			)
			return
		}

		editor := module.Editor{}
		err = editor.Edit(absPath)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		log.WithFields(log.Fields{
			"file": absPath,
		}).Debug("Service file created")

		fmt.Println(Green("Service file created successfully"))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
