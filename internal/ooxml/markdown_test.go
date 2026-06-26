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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInsertMarkdownTable verifies that a GFM pipe table in Markdown input is
// converted into an OOXML <w:tbl> element.  This test fails without the
// goldmark table extension and passes once it is enabled.
func TestInsertMarkdownTable(t *testing.T) {
	a := assert.New(t)

	md := `
| Name    | IP        |
|---------|-----------|
| switch1 | 10.0.0.1  |
| switch2 | 10.0.0.2  |
`
	doc := newTestDocument()
	a.NoError(doc.InsertMarkdown(strings.NewReader(md), "main"))

	// The body must contain a table element.
	body := doc.document.doc.FindElement("//w:body")
	a.NotNil(body, "body element not found")
	tbl := body.FindElement("//w:tbl")
	a.NotNil(tbl, "expected a w:tbl element from markdown table, got none")
}
