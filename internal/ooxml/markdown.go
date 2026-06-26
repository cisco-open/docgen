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
package ooxml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// InsertMarkdown converts Markdown content to HTML and inserts it into the document
// at the specified bookmark location.
//
// This method uses goldmark to convert Markdown to HTML, then delegates to InsertHTML.
//
// Returns an error if the bookmark is not found or if there are issues parsing
// the Markdown or inserting the HTML.
func (doc *Document) InsertMarkdown(input io.Reader, bookmark string) error {
	markdownBytes, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read markdown: %w", err)
	}

	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	var htmlBuf bytes.Buffer
	if err := md.Convert(markdownBytes, &htmlBuf); err != nil {
		return fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	return doc.InsertHTML(&htmlBuf, bookmark)
}
