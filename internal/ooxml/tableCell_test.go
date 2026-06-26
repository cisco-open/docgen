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

// TestColSpan validates colspan functionality
func TestColspan(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<th colspan="2">a1</th>
			<th>c1</th>
		</tr>
		<tr>
			<td>a2</td>
			<td>b2</td>
			<td>c2</td>
		</tr>
	</table>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "colspan.docx"))

	// Write to test doc
	doc = newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	r := func(txt string, props ...prop) *etree.Element {
		tc := newTableCellEl(Props(props...)).e
		tc.AddChild(newPara(newRun(newText(txt))))
		return tc
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newTable(
				newTableGrid(),
				newTableRow(
					headerCell("a1", tableCellColspan{2}),
					headerCell("c1"),
				),
				newTableRow(r("a2"), r("b2"), r("c2")),
			),
		),
	)

	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestRowSpan validates rowspan functionality
func TestRowspan(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<th>a1</th>
			<th>b1</th>
		</tr>
		<tr>
			<td rowspan="2">a2</td>
			<td>b2</td>
		</tr>
		<tr>
			<td>b3</td>
		</tr>
		<tr>
			<td rowspan="2">a4</td>
			<td>b4</td>
		</tr>
		<tr>
			<td>b5</td>
		</tr>
	</table>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "rowspan.docx"))

	// Write to test doc
	doc = newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	r := func(txt string, props ...prop) *etree.Element {
		tc := newTableCellEl(Props(props...)).e
		tc.AddChild(newPara(newRun(newText(txt))))
		return tc
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newTable(
				newTableGrid(),
				newTableRow(headerCell("a1"), headerCell("b1")),
				newTableRow(
					r("a2", tableCellRowspan{true}),
					r("b2"),
				),
				newTableRow(
					newTableCellEl(Props(tableCellRowspan{})).
						AddXML(newPara()).e,
					r("b3"),
				),
				newTableRow(
					r("a4", tableCellRowspan{true}),
					r("b4"),
				),
				newTableRow(
					newTableCellEl(Props(tableCellRowspan{})).
						AddXML(newPara()).e,
					r("b5"),
				),
			),
		),
	)
	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestRowColSpan validates rowspan and colspan functionality
func TestRowColspan(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<th>a1</th>
			<th>b1</th>
			<th>c1</th>
		</tr>
		<tr>
			<td rowspan="2" colspan="2">a2</td>
			<td>c2</td>
		</tr>
		<tr>
			<td>c3</td>
		</tr>
	</table>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "rowcolspan.docx"))

	// Write to test doc
	doc = newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	r := func(txt string, props ...prop) *etree.Element {
		tc := newTableCellEl(Props(props...)).e
		tc.AddChild(newPara(newRun(newText(txt))))
		return tc
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newTable(
				newTableGrid(),
				newTableRow(headerCell("a1"), headerCell("b1"), headerCell("c1")),
				newTableRow(r("a2", tableCellColspan{2}, tableCellRowspan{true}), r("c2")),
				newTableRow(
					newTableCellEl(Props(tableCellColspan{2}, tableCellRowspan{})).
						AddXML(newPara()).e,
					r("c3"),
				),
			),
		),
	)
	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}
