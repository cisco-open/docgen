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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	cfg := New()
	a.Equal("", cfg.Context)
	a.Equal("input.html", cfg.Input)
	a.Equal("out.docx", cfg.Output)
	a.Equal("main", cfg.Bookmark)
	a.Equal("docgen.log", cfg.Logfile)
	a.False(cfg.Verbose)
	a.Equal("", cfg.Workdir)
}

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Valid config
	validConfig := `
context: test-context.json
input: test-input.md
output: test-output.docx
template: test-template.docx
bookmark: test-bookmark
logfile: test-logfile.log
verbose: true
workdir: /tmp/test
`
	err := os.WriteFile(configPath, []byte(validConfig), 0644)
	a.NoError(err)

	cfg, err := LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)
	a.Equal("test-context.json", filepath.Base(cfg.Context))
	a.Equal("test-input.md", filepath.Base(cfg.Input))
	a.Equal("test-output.docx", filepath.Base(cfg.Output))
	a.Equal("test-template.docx", filepath.Base(cfg.Template))
	a.Equal("test-bookmark", cfg.Bookmark)
	a.Equal("test-logfile.log", filepath.Base(cfg.Logfile))
	a.True(cfg.Verbose)
}

func TestLoadConfigPartial(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Partial config - should use defaults for missing values
	partialConfig := `
input: custom-input.html
output: custom-output.docx
`
	err := os.WriteFile(configPath, []byte(partialConfig), 0644)
	a.NoError(err)

	cfg, err := LoadConfig(configPath)
	a.NoError(err)
	a.NotNil(cfg)
	a.Equal("custom-input.html", filepath.Base(cfg.Input))
	a.Equal("custom-output.docx", filepath.Base(cfg.Output))
	// Should have defaults for missing fields (context is now optional, defaults to empty)
	a.Equal("", cfg.Context)
	a.Equal("main", cfg.Bookmark)
	a.Equal("docgen.log", filepath.Base(cfg.Logfile))
	a.False(cfg.Verbose)
}

func TestLoadConfigNonExistent(t *testing.T) {
	a := assert.New(t)

	cfg, err := LoadConfig("/nonexistent/config.yaml")
	a.Error(err)
	a.Nil(cfg)
	a.Contains(err.Error(), "failed to read config file")
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `
input: test.html
output: [this is not valid yaml
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	a.NoError(err)

	cfg, err := LoadConfig(configPath)
	a.Error(err)
	a.Nil(cfg)
	a.Contains(err.Error(), "failed to parse config file")
}

func TestApplyDefaults(t *testing.T) {
	a := assert.New(t)

	cfg := Config{
		Input:  "test.html",
		Output: "test.docx",
		// Other fields empty
	}

	cfg.ApplyDefaults()

	// Should have defaults for empty fields (context is now optional, defaults to empty)
	a.Equal("", cfg.Context)
	a.Equal("main", cfg.Bookmark)
	a.Equal("docgen.log", cfg.Logfile)
	// Should keep custom values
	a.Equal("test.html", cfg.Input)
	a.Equal("test.docx", cfg.Output)
}

func TestNormalizePaths(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()

	// Create some test files in tmpDir
	contextFile := filepath.Join(tmpDir, "context.json")
	err := os.WriteFile(contextFile, []byte("{}"), 0644)
	a.NoError(err)

	cfg := Config{
		Context: contextFile,
		Input:   "input.html",
		Output:  "output.docx",
		Workdir: tmpDir,
	}

	cfg.NormalizePaths()

	// Context should be absolute (exists in filesystem)
	a.True(filepath.IsAbs(cfg.Context))
	// Input and Output should be relative to workdir (don't exist)
	a.Contains(cfg.Input, tmpDir)
	a.Contains(cfg.Output, tmpDir)
}

func TestNormalizePathsEmptyStrings(t *testing.T) {
	a := assert.New(t)

	cfg := Config{
		Context:  "",
		Input:    "",
		Output:   "",
		Template: "",
		Logfile:  "",
	}

	cfg.NormalizePaths()

	// Empty strings should remain empty
	a.Equal("", cfg.Context)
	a.Equal("", cfg.Input)
	a.Equal("", cfg.Output)
	a.Equal("", cfg.Template)
	a.Equal("", cfg.Logfile)
}

func TestNormalizePathsWithSlashes(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()

	cfg := Config{
		Context: "path/to/context.json",
		Input:   "path/to/input.html",
		Workdir: tmpDir,
	}

	cfg.NormalizePaths()

	// Paths should be converted to OS-specific separators and joined with workdir
	a.Contains(cfg.Context, tmpDir)
	a.Contains(cfg.Input, tmpDir)
	// Paths should use filepath separator (not forward slash)
	a.Equal(filepath.Join(tmpDir, "path", "to", "context.json"), cfg.Context)
	a.Equal(filepath.Join(tmpDir, "path", "to", "input.html"), cfg.Input)
}

func TestNormalizePathsExistingVsNonExisting(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()

	// Create an existing file in tmpDir
	existingFile := filepath.Join(tmpDir, "existing.html")
	err := os.WriteFile(existingFile, []byte("<p>test</p>"), 0644)
	a.NoError(err)

	// Save original working directory and change to tmpDir
	originalWd, err := os.Getwd()
	a.NoError(err)
	err = os.Chdir(tmpDir)
	a.NoError(err)
	defer os.Chdir(originalWd)

	workdir := filepath.Join(tmpDir, "other-workdir")
	err = os.MkdirAll(workdir, 0755)
	a.NoError(err)

	cfg := Config{
		Input:   "existing.html",    // exists in CWD
		Output:  "nonexisting.docx", // doesn't exist, should use workdir
		Workdir: workdir,
	}

	cfg.NormalizePaths()

	// Existing file should be cleaned but not joined with workdir
	// (it exists in CWD, so os.Stat succeeds and we don't join with workdir)
	a.Equal("existing.html", cfg.Input)

	// Non-existing file should be joined with workdir
	a.Contains(cfg.Output, workdir)
	a.Contains(cfg.Output, "nonexisting.docx")
}
