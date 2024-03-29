// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"path/filepath"
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

// From var
var From string

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Interact with one of the configured services",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Call command got called.")

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

		files := make(map[string]util.File)

		if From == "" || !util.FileExists(From) {
			files, err = util.ListFiles(util.EnsureTrailingSlash(conf.Services.Directory))
		} else {
			files[filepath.Base(From)] = util.File{
				Path: From,
				Name: filepath.Base(From),
			}
		}

		if err != nil {
			fmt.Printf(
				"Error while listing services under %s: %s",
				util.EnsureTrailingSlash(conf.Services.Directory),
				err.Error(),
			)
			return
		}

		data := []string{}
		index := map[string]*model.Service{}
		service := &model.Service{}

		for _, v := range files {
			if strings.Contains(v.Name, ".toml") {

				service = model.NewEmptyService(v.Name)
				err = service.Decode(v.Path)

				if err != nil {
					fmt.Printf(
						"Error while decoding service %s: %s",
						v.Path,
						err.Error(),
					)
					return
				}

				for _, end := range service.Endpoint {
					data = append(
						data,
						fmt.Sprintf("%s - %s", service.Main.ID, end.ID),
					)

					index[fmt.Sprintf("%s - %s", service.Main.ID, end.ID)] = service
				}
			}
		}

		result := ""
		finder := module.FuzzyFinder{}
		prompt := module.Prompt{}

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

		caller := module.NewCaller(module.NewHTTPClient())
		fields := caller.GetFields(result, index[result])

		val := ""

		for key, field := range fields {
			if field.IsOptional {
				val, err = prompt.Input(
					field.Prompt,
					module.Optional,
				)

				if module.IsEmpty(val) {
					val = fields[key].Default
				}

			} else {
				val, err = prompt.Input(
					field.Prompt,
					module.NotEmpty,
				)
			}

			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				return
			}

			fields[key] = module.Field{
				Prompt:     field.Prompt,
				IsOptional: field.IsOptional,
				Default:    field.Default,
				Value:      val,
			}
		}

		spin := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
		spin.Color("green")
		spin.Start()

		response, err := caller.Call(result, index[result], fields)

		spin.Stop()

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		if response == nil {
			fmt.Println(Red("Invalid Response!"))
			return
		}

		fmt.Println(caller.Pretty(response))
	},
}

func init() {
	callCmd.PersistentFlags().StringVarP(
		&From,
		"from",
		"f",
		"./.poodle.toml",
		"service definition file",
	)
}

func init() {
	rootCmd.AddCommand(callCmd)
}
