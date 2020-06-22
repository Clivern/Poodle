// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync API services definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`WIP`)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
