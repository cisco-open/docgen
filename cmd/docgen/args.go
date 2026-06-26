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
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/ciscotools/docgen/pkg/config"
)

// Populated from CI/CD
var version = "(dev)"

// Args are CLI arguments.
type Args struct {
	Config   string `arg:"--config"                     help:"Path to YAML configuration file"`
	CAF      bool   `arg:"--caf"                        help:"Run in CAF service mode"`
	Context  string `arg:"-c,--context,env:CONTEXT"     help:"Tag context JSON file (optional)"`
	Input    string `arg:"-i,--input,env:INPUT"         help:"Input file (.html, .md, or .markdown)"          default:"input.html"`
	Output   string `arg:"-o,--output"                  help:"Output file"                                     default:"out.docx"`
	Template string `arg:"--template"                   help:"Use a custom docx template"`
	Bookmark string `arg:"--bookmark"                   help:"Bookmark to render into"                         default:"main"`
	Logfile  string `arg:"--logfile"                    help:"Logfile"                                         default:"docgen.log"`
	Verbose  bool   `arg:"-v,--verbose"                 help:"Verbose output"`
	Workdir  string `arg:"--workdir"                    help:"Working directory"`
}

// Description is the CLI description string.
func (Args) Description() string {
	return "Docgen Word document generator"
}

// Version is the CLI version string, populated by goreleaser.
func (Args) Version() string {
	return fmt.Sprintf("docgen %s", version)
}

// readArgs collects the CLI args and returns a config.Config.
// If --caf is provided, configuration is read from /share/job.json.
// If a config file is provided, CLI parameters are ignored.
func readArgs() config.Config {
	var args Args
	arg.MustParse(&args)

	// CAF mode: read job.json from /share
	if args.CAF {
		cfg, err := readCAF()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CAF job: %v\n", err)
			os.Exit(1)
		}
		return cfg
	}

	// If config file is provided, load from it (ignoring CLI params)
	if args.Config != "" {
		cfg, err := config.LoadConfig(args.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
			os.Exit(1)
		}
		return *cfg
	}

	// Build config from CLI arguments
	cfg := config.New()
	cfg.Context = args.Context
	cfg.Input = args.Input
	cfg.Output = args.Output
	cfg.Template = args.Template
	cfg.Bookmark = args.Bookmark
	cfg.Logfile = args.Logfile
	cfg.Verbose = args.Verbose
	cfg.Workdir = args.Workdir

	// Apply path normalization using the config package method
	cfg.NormalizePaths()

	return cfg
}
