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
package docgen_test

import (
	"strings"
	"testing"

	"github.com/ciscotools/docgen"
	"github.com/ciscotools/docgen/internal/templates"
	"github.com/stretchr/testify/assert"
)

// TestDocumentInsertMarkdown verifies that Document.InsertMarkdown can be called
func TestDocumentInsertMarkdown(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	markdown := strings.NewReader("# Test Heading\n\nThis is **bold** text and this is *italic* text.\n\n- Item 1\n- Item 2\n- Item 3")
	err = doc.InsertMarkdown(markdown, "main")
	if err != nil {
		t.Fatalf("InsertMarkdown failed: %v", err)
	}
}

// TestDocumentInsertMarkdownWithTable verifies that Markdown tables can be inserted
func TestDocumentInsertMarkdownWithTable(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	markdown := strings.NewReader(`
# Table Example

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
| Data 4   | Data 5   | Data 6   |
`)
	err = doc.InsertMarkdown(markdown, "main")
	if err != nil {
		t.Fatalf("InsertMarkdown with table failed: %v", err)
	}
}

// TestDocumentInsertMarkdownWithLinks verifies that Markdown links can be inserted
func TestDocumentInsertMarkdownWithLinks(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	markdown := strings.NewReader("Check out [this link](https://example.com) for more info.")
	err = doc.InsertMarkdown(markdown, "main")
	if err != nil {
		t.Fatalf("InsertMarkdown with links failed: %v", err)
	}
}

// TestDocumentInsertMarkdownInvalidBookmark verifies error handling for invalid bookmarks
func TestDocumentInsertMarkdownInvalidBookmark(t *testing.T) {
	a := assert.New(t)
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	markdown := strings.NewReader("# Test")
	err = doc.InsertMarkdown(markdown, "nonexistent")
	a.Error(err, "Expected error for nonexistent bookmark")
	a.Contains(err.Error(), "nonexistent", "Error should mention the bookmark name")
}
