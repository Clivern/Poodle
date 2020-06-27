// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
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

		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		}

		gh_username := promptui.Prompt{
			Label:     "Github username",
			Templates: templates,
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("Invalid input")
				}
				return nil
			},
		}

		result, err := gh_username.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You answered %s\n", result)

		fmt.Println(`WIP`)
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
