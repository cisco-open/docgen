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

func TestToC(t *testing.T) {
	a := assert.New(t)

	html := `
	<h1>Header 1</h1>
	<h2>Header 1.1</h2>
	<h3>Header 1.1.1</h3>
	<h3>Header 1.1.2</h3>
	<h1>Header 2</h1>
	<div class="toc" />
	`

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "toc.docx"))
}
