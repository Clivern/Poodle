// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package config

// Config struct
type Config struct {
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

// RequestItem request a value from end user
func (c *Config) RequestItem(label string) (string, error) {
	return "", nil
}

// UpdateItem updates an item
func (c *Config) UpdateItem(key, value string) (bool, error) {
	return true, nil
}
