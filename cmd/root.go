// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Verbose var
var Verbose bool

var rootCmd = &cobra.Command{
	Use: "poodle",
	Short: `Work seamlessly with Poodle from the command line.

Poodle is in early stages of development, and we'd love to hear your
feedback at <https://github.com/Clivern/Poodle>`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

// Execute runs cmd tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
