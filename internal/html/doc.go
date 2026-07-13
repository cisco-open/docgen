// Package html performs HTML templating

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
package html

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cisco-open/docgen/internal/log"

	"github.com/Masterminds/sprig/v3"
)

// Template is an html/template wrapper
type Template struct {
	template *template.Template
}

// Render renders a template
func (t *Template) Render(ctx any) (string, error) {
	var b strings.Builder
	err := t.template.Execute(&b, ctx)
	return b.String(), err
}

// Document is a new HTML document
type Document struct {
	ctx             any
	entry           string
	Templates       map[string]*Template
	TemplateOptions []string
}

// New creates a new document
func New(entry string, ctx any) *Document {
	return &Document{
		ctx:       ctx,
		entry:     entry,
		Templates: map[string]*Template{},
	}
}

// RenderTemplates pre-renders all templates
func (doc *Document) RenderTemplates() {
	var wg sync.WaitGroup
	for name, tpl := range doc.Templates {
		wg.Add(1)
		name, tpl := name, tpl
		go func() {
			if _, err := tpl.Render(doc.ctx); err != nil {
				log.Warn().
					Err(err).
					Str("name", name).
					Msg("error rendering template")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// NewTemplate creates a new template without registering it
func (doc *Document) NewTemplate(
	name string,
	r io.Reader,
	funcs ...template.FuncMap,
) (*template.Template, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read body template body %w", err)
	}
	t := template.New(name).
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"debug":   debugHelper,
			"date":    dateHelper,
			"include": doc.includeHelper,
			"R":       doc.rootHelper,
			"root":    doc.rootHelper,
		})
	for _, fs := range funcs {
		t.Funcs(fs)
	}
	for _, opt := range doc.TemplateOptions {
		t = t.Option(opt)
	}
	t, err = t.Parse(string(body))
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s %w", name, err)
	}
	return t, nil
}

// RegisterTemplate registers a single template
func (doc *Document) RegisterTemplate(name string, r io.Reader, funcs ...template.FuncMap) error {
	t, err := doc.NewTemplate(name, r, funcs...)
	if err != nil {
		return err
	}
	doc.Templates[name] = &Template{t}
	return nil
}

// RegisterTemplatePath registers all templates in a folder
func (doc *Document) RegisterTemplatePath(path string) error {
	log.Debug().Msgf("registering templates in %s", path)
	return doc.RegisterTemplates(os.DirFS(path))
}

// RegisterTemplates registers templates in a fs.FS filesystem
func (doc *Document) RegisterTemplates(fileSystem fs.FS) error {
	return fs.WalkDir(fileSystem, ".",
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				log.Error().Err(err).Msg("error reading templates")
				return nil
			}
			// Recursively follow symlinks
			if info.Type()&fs.ModeSymlink != 0 {
				linkedPath, err := os.Readlink(path)
				if err != nil {
					return nil
				}
				doc.RegisterTemplates(os.DirFS(linkedPath))
			}
			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(info.Name())
			if ext != ".html" {
				return nil
			}
			name := strings.TrimSuffix(info.Name(), ext)
			if name == "template" {
				name = filepath.Base(filepath.Dir(path))
			}
			if _, ok := doc.Templates[name]; ok {
				log.Warn().
					Str("name", name).
					Msg("template already registered")
				return fs.SkipDir
			}
			r, err := fileSystem.Open(path)
			if err != nil {
				log.Error().Err(err).Msgf("cannot read file %s", path)
				return nil
			}
			defer r.Close()
			log.Debug().Msgf("registering template %s", name)
			if err := doc.RegisterTemplate(name, r); err != nil {
				log.Error().Err(err).Msgf("cannot register template %s", name)
			}
			return nil
		})
}

// RenderHTML renders the HTML document
func (doc *Document) RenderHTML(w io.Writer) error {
	tpl, ok := doc.Templates[doc.entry]
	if !ok {
		return fmt.Errorf("entry %s not registered", doc.entry)
	}
	return tpl.template.Execute(w, doc.ctx)
}
