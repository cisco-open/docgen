// Package config provides YAML configuration file support for docgen.

// SPDX-License-Identifier: Apache-2.0

// Copyright 2026 Cisco Systems, Inc. and their affiliates

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for docgen.
type Config struct {
	Context  string `yaml:"context"`
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Template string `yaml:"template"`
	Bookmark string `yaml:"bookmark"`
	Logfile  string `yaml:"logfile"`
	Verbose  bool   `yaml:"verbose"`
	Workdir  string `yaml:"workdir"`
}

// New returns a Config with default values.
func New() Config {
	return Config{
		Context:  "",
		Input:    "input.html",
		Output:   "out.docx",
		Bookmark: "main",
		Logfile:  "docgen.log",
		Verbose:  false,
		Workdir:  "",
	}
}

// LoadConfig reads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.ApplyDefaults()
	cfg.NormalizePaths()

	return &cfg, nil
}

// ApplyDefaults sets default values for missing fields.
func (c *Config) ApplyDefaults() {
	defaults := New()
	if c.Input == "" {
		c.Input = defaults.Input
	}
	if c.Output == "" {
		c.Output = defaults.Output
	}
	if c.Bookmark == "" {
		c.Bookmark = defaults.Bookmark
	}
	if c.Logfile == "" {
		c.Logfile = defaults.Logfile
	}
}

// NormalizePaths normalizes file paths relative to workdir.
func (c *Config) NormalizePaths() {
	if c.Workdir == "" {
		if wd, err := os.Getwd(); err == nil {
			c.Workdir = wd
		}
	} else {
		_ = os.MkdirAll(c.Workdir, 0o750)
	}

	clean := func(s string) string {
		if s == "" {
			return s
		}
		s = filepath.FromSlash(s)
		// Only use workdir if the path is not present in the CWD
		if _, err := os.Stat(s); os.IsNotExist(err) {
			s = filepath.Join(c.Workdir, s)
		}
		s = filepath.Clean(s)
		return s
	}

	c.Input = clean(c.Input)
	c.Output = clean(c.Output)
	c.Context = clean(c.Context)
	c.Logfile = clean(c.Logfile)
	c.Template = clean(c.Template)
}
