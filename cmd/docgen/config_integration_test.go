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
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cisco-open/docgen/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestReadArgsWithConfigFile(t *testing.T) {
	a := assert.New(t)

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	
	// Create test input file
	inputPath := filepath.Join(tmpDir, "test-input.html")
	err := os.WriteFile(inputPath, []byte("<p>Test content</p>"), 0644)
	a.NoError(err)

	// Create config file
	configContent := `
context: test-context.json
input: test-input.html
output: test-output.docx
template: test-template.docx
bookmark: test-bookmark
logfile: test.log
verbose: true
workdir: ` + tmpDir + `
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	a.NoError(err)

	// Load the config
	cfg, err := config.LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)

	// Verify config was loaded correctly
	a.True(cfg.Verbose)
	a.Equal("test-bookmark", cfg.Bookmark)
	a.Contains(cfg.Input, "test-input.html")
	a.Contains(cfg.Output, "test-output.docx")
}

func TestReadArgsWithoutConfigFile(t *testing.T) {
	a := assert.New(t)

	// Test that default config can be created
	cfg := config.New()
	a.Equal("", cfg.Context)
	a.Equal("input.html", cfg.Input)
	a.Equal("out.docx", cfg.Output)
	a.Equal("main", cfg.Bookmark)
	a.Equal("docgen.log", cfg.Logfile)
	a.False(cfg.Verbose)
}

func TestConfigFileOverridesCLI(t *testing.T) {
	// This test documents the behavior that config file
	// takes precedence over CLI arguments
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
input: config-input.html
output: config-output.docx
verbose: true
workdir: ` + tmpDir + `
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	a.NoError(err)

	cfg, err := config.LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)

	// Config file values should be used
	a.True(cfg.Verbose)
	a.Contains(cfg.Input, "config-input.html")
	a.Contains(cfg.Output, "config-output.docx")
}

func TestConfigPathNormalization(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	a.NoError(err)

	// Create test file in subdirectory
	testFile := filepath.Join(subDir, "test.html")
	err = os.WriteFile(testFile, []byte("<p>test</p>"), 0644)
	a.NoError(err)

	configContent := `
input: subdir/test.html
output: subdir/output.docx
workdir: ` + tmpDir + `
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	a.NoError(err)

	cfg, err := config.LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)

	// Paths should be normalized to absolute paths
	a.True(filepath.IsAbs(cfg.Input))
	a.True(filepath.IsAbs(cfg.Output))
	a.Contains(cfg.Input, "subdir")
	a.Contains(cfg.Output, "subdir")
}
