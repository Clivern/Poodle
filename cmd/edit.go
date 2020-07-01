// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"

	"github.com/clivern/poodle/core/model"
	"github.com/clivern/poodle/core/module"
	"github.com/clivern/poodle/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit service definition file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Edit command got called.")

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
		finder := module.FuzzyFinder{}

		localFiles, err := util.ListFiles(util.EnsureTrailingSlash(conf.Services.Directory))

		if err != nil {
			fmt.Printf(
				"Error while listing local files in %s: %s",
				util.EnsureTrailingSlash(conf.Services.Directory),
				err.Error(),
			)
			return
		}

		files := []string{}

		for _, file := range localFiles {
			if !strings.Contains(file.Name, ".toml") {
				continue
			}

			files = append(files, strings.Replace(file.Name, ".toml", "", -1))
		}

		relPath := ""

		if finder.Available() {
			relPath, err = finder.Show(files)
		} else {
			relPath, err = prompt.Select(
				fmt.Sprintf("Select a Service"),
				files,
			)
		}

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		absPath := fmt.Sprintf(
			"%s%s.toml",
			util.EnsureTrailingSlash(conf.Services.Directory),
			relPath,
		)

		editor := module.Editor{}
		err = editor.Edit(absPath)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		log.WithFields(log.Fields{
			"file": absPath,
		}).Debug("Service file updated")

		fmt.Println("Service file updated")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
