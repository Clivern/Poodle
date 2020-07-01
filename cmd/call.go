// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/clivern/poodle/core/module"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Interact with one of the configured API services",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Call command got called.")

		finder := module.FuzzyFinder{}
		prompt := module.Prompt{}

		data := []string{
			"A",
			"B",
			"C",
			"D",
			"E",
			"F",
			"G",
			"H",
			"I",
		}

		result := ""

		if finder.Available() {
			result, err = finder.Show(data)
		} else {
			result, err = prompt.Select(
				fmt.Sprintf("Select an Endpoint"),
				data,
			)
		}

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(callCmd)
}
