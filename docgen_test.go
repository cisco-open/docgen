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
	"bytes"
	"strings"
	"testing"

	"github.com/cisco-open/docgen"
	"github.com/cisco-open/docgen/internal/templates"
	"github.com/stretchr/testify/assert"
)

// TestNewDocument verifies that NewDocument can be called from external packages
func TestNewDocument(t *testing.T) {
	a := assert.New(t)
	_, err := docgen.NewDocument(templates.MainTemplate)
	a.NoError(err)
}

// TestDocumentWrite verifies that Document.Write can be called
func TestDocumentWrite(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	var buf bytes.Buffer
	err = doc.Write(&buf)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Write produced empty output")
	}
}

// TestDocumentInsertHTML verifies that Document.InsertHTML can be called
func TestDocumentInsertHTML(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	html := strings.NewReader("<p>Test content</p>")
	err = doc.InsertHTML(html, "main")
	if err != nil {
		t.Fatalf("InsertHTML failed: %v", err)
	}
}

// TestDocumentRenderTags verifies that Document.RenderTags can be called
func TestDocumentRenderTags(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	ctx := map[string]any{
		"title": "Test Document",
	}
	err = doc.RenderTags(ctx)
	if err != nil {
		t.Fatalf("RenderTags failed: %v", err)
	}
}

// TestDocumentGetBookmark verifies that Document.GetBookmark can be called
func TestDocumentGetBookmark(t *testing.T) {
	doc, err := docgen.NewDocument(templates.MainTemplate)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}

	bookmark := doc.GetBookmark("main")
	if bookmark == nil {
		t.Fatal("GetBookmark returned nil for 'main' bookmark")
	}
}
