// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/clivern/poodle/core/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Poodle",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Configure command got called.")

		conf := config.Config{}

		configureWith, err := conf.Select(
			fmt.Sprintf("Configure With"),
			[]string{"Interactive", fmt.Sprintf("Editor (%s)", os.Getenv("EDITOR"))},
		)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		githubUsername, err := conf.Prompt("Github Username:", config.NotEmpty)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		githubToken, err := conf.Prompt("Github OAuth Token:", config.NotEmpty)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		fmt.Println(configureWith)
		fmt.Println(githubUsername)
		fmt.Println(githubToken)
		fmt.Println(`WIP`)
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
