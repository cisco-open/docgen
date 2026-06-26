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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ciscotools/docgen/pkg/config"
)

const cafSharePath = "/share"

func getCAFSharePath() string {
	path := cafSharePath
	if _, err := os.Stat(path); err != nil {
		path, _ = os.Getwd()
	}
	return path
}

// cafJob represents the CAF job.json structure.
type cafJob struct {
	ID       string `json:"id"`
	Metadata struct {
		Customer string         `json:"customer"`
		Title    string         `json:"title"`
		Bookmark string         `json:"bookmark"`
		Extra    map[string]any `json:"extra"`
	} `json:"metadata"`
	Files []struct {
		Name string `json:"name"`
	} `json:"files"`
}

// toConfig converts a cafJob to a docgen config.Config.
// It writes job metadata to /share/context.json for use as template context.
func (job cafJob) toConfig() (config.Config, error) {
	if len(job.Files) == 0 {
		return config.Config{}, fmt.Errorf("no files provided in job")
	}

	sharePath := getCAFSharePath()

	// Build template context from job metadata and write to context.json
	ctx := map[string]any{
		"customer": job.Metadata.Customer,
		"title":    job.Metadata.Title,
		"jobID":    job.ID,
	}
	for k, v := range job.Metadata.Extra {
		ctx[k] = v
	}
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to marshal context: %w", err)
	}
	ctxPath := filepath.Join(sharePath, "context.json")
	if err := os.WriteFile(ctxPath, ctxBytes, 0o644); err != nil {
		return config.Config{}, fmt.Errorf("failed to write context.json: %w", err)
	}

	cfg := config.New()
	cfg.Input = filepath.Join(sharePath, job.Files[0].Name)
	cfg.Output = filepath.Join(sharePath, "output.docx")
	cfg.Context = ctxPath
	cfg.Workdir = sharePath
	cfg.Verbose = true

	if job.Metadata.Bookmark != "" {
		cfg.Bookmark = job.Metadata.Bookmark
	}

	return cfg, nil
}

// readCAF reads and parses /share/job.json and returns a config.Config.
func readCAF() (config.Config, error) {
	path := filepath.Join(getCAFSharePath(), "job.json")

	var job cafJob
	job.Metadata.Extra = make(map[string]any)

	b, err := os.ReadFile(path)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to read job.json: %w", err)
	}
	if err := json.Unmarshal(b, &job); err != nil {
		return config.Config{}, fmt.Errorf("failed to parse job.json: %w", err)
	}

	return job.toConfig()
}
