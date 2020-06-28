// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Config struct
type Config struct {
}

// NotEmpty returns error if input is empty
func NotEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("Input must not be empty")
	}
	return nil
}

// Exists check if file exists
func (c *Config) Exists() (bool, error) {
	return true, nil
}

// Create creates a config file
func (c *Config) Create() (bool, error) {
	return true, nil
}

// IsItemEmpty checks if config item not empty
func (c *Config) IsItemEmpty(key string) (bool, error) {
	return true, nil
}

// Prompt request a value from end user
func (c *Config) Prompt(label string, validate promptui.ValidateFunc) (string, error) {

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	item := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := item.Run()

	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	return result, nil
}

// UpdateItem updates an item
func (c *Config) UpdateItem(key, value string) (bool, error) {
	return true, nil
}
