// Package main implements the docgen command line tool.

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
	"os"
	"path/filepath"
	"strings"

	"github.com/cisco-open/docgen/internal/ooxml"
	"github.com/cisco-open/docgen/internal/templates"
	"github.com/cockroachdb/errors"
)

func loadTemplate(templateFile string) (*ooxml.Document, error) {

	if _, err := os.Stat(templateFile); err == nil {
		body, err := os.ReadFile(templateFile)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return ooxml.NewDocument(body)
	} else {
		return ooxml.NewDocument(templates.MainTemplate)
	}
}

func renderContext(ctxFile string, doc *ooxml.Document) error {
	if ctxFile == "" {
		// No context provided; render tags with an empty context so that any
		// {{tag}} placeholders in the template are replaced with empty strings.
		return errors.WithStack(doc.RenderTags(map[string]any{}))
	}
	ctxBytes, err := os.ReadFile(ctxFile)
	if err != nil {
		return errors.WithStack(err)
	}
	var ctx any
	if err := json.Unmarshal(ctxBytes, &ctx); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(doc.RenderTags(ctx))
}

func renderHTML(htmlFile string, bookmark string, doc *ooxml.Document) error {
	if htmlFile == "" {
		return nil
	}
	f, err := os.Open(htmlFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	return errors.WithStack(doc.InsertHTML(f, bookmark))
}

func renderMarkdown(markdownFile string, bookmark string, doc *ooxml.Document) error {
	if markdownFile == "" {
		return nil
	}
	f, err := os.Open(markdownFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	return errors.WithStack(doc.InsertMarkdown(f, bookmark))
}

func renderInput(inputFile string, bookmark string, doc *ooxml.Document) error {
	if inputFile == "" {
		return nil
	}

	// Determine file type by extension
	ext := strings.ToLower(filepath.Ext(inputFile))
	switch ext {
	case ".md", ".markdown":
		return renderMarkdown(inputFile, bookmark, doc)
	case ".html", ".htm":
		return renderHTML(inputFile, bookmark, doc)
	default:
		// Default to HTML for backwards compatibility
		return renderHTML(inputFile, bookmark, doc)
	}
}

func main() {
	cfg := readArgs()
	doc, err := loadTemplate(cfg.Template)
	if err != nil {
		panic("cannot create Word document")
	}
	if err := renderContext(cfg.Context, doc); err != nil {
		panic("cannot render context")
	}
	if err := renderInput(cfg.Input, cfg.Bookmark, doc); err != nil {
		panic("cannot render input")
	}

	if err := doc.WriteFile(cfg.Output); err != nil {
		panic("cannot write Word document")
	}
}
