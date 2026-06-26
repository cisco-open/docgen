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
	"strconv"
	"strings"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	a := assert.New(t)

	html := `
	<ol>
	  <li>one</li>
	  <li>two</li>
	  <ol>
	    <li>three</li>
	  </ol>
	</ol>
	<ul>
	  <li>a</li>
	  <li>b</li>
	<ul>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "list.docx"))

	// Write to test doc
	doc = newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	item := func(txt string, lvl, id int) *etree.Element {
		return newPara(
			tag("w:pPr",
				tag("w:pStyle", attr("w:val", "ListParagraph")),
				tag("w:numPr",
					tag("w:ilvl", attr("w:val", strconv.Itoa(lvl))),
					tag("w:numId", attr("w:val", strconv.Itoa(id))),
				),
			),
			newRun(newText(txt)))
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			item("one", 0, 1),
			item("two", 0, 1),
			item("three", 1, 1),
			item("a", 0, 2),
			item("b", 0, 2),
		),
	)
	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}
