// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync API services definitions",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Sync command got called.")

		fmt.Println(`WIP`)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
