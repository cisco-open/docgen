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

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestTableCaption(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<caption>test</caption>
		<tr>
			<th>header1</th>
			<th>header2</th>
		</tr>
		<tr>
			<td>cell1</td>
			<td>cell2</td>
		</tr>
	</table>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "tableCaption.docx"))

	// Write to test doc
	doc = newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	r := func(txts ...string) *etree.Element {
		tr := newTableRow()
		for _, txt := range txts {
			tr.AddChild(newTableCell(newPara(newRun(newText(txt)))))
		}
		return tr
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newPara(
				tag("w:pPr",
					tag("w:pStyle", attr("w:val", "Caption")),
				),
				newRun(newText("test")),
			),
			newTable(
				newTableGrid(),
				newTableRow(headerCell("header1"), headerCell("header2")),
				r("cell1", "cell2"),
			),
		),
	)
	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}
