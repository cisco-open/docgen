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

// headerCell builds the expected OOXML for a th cell: navy background, bold white text.
func headerCell(txt string, extraProps ...prop) *etree.Element {
	tc := newTableCellEl(Props(append([]prop{tableCellHeader{}}, extraProps...)...)).e
	tc.AddChild(newPara(newRun(
		tag("w:rPr",
			tag("w:b"),
			tag("w:color", attr("w:val", "FFFFFF"), attr("w:themeColor", "background2")),
		),
		newText(txt),
	)))
	return tc
}

// subheaderCell builds the expected OOXML for a th/td.subheader cell: cyan background, bold white text.
func subheaderCell(txt string, extraProps ...prop) *etree.Element {
	tc := newTableCellEl(Props(append([]prop{tableCellSubheader{}}, extraProps...)...)).e
	tc.AddChild(newPara(newRun(
		tag("w:rPr",
			tag("w:b"),
			tag("w:color", attr("w:val", "FFFFFF"), attr("w:themeColor", "background2")),
		),
		newText(txt),
	)))
	return tc
}

func TestTable(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<th>a1</th>
			<th>b1</th>
		</tr>
		<tr>
			<td>a2</td>
			<td>b2</td>
		</tr>
		<tr>
			<td>a3</td>
			<td>b3</td>
		</tr>
	</table>
  `

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "table.docx"))

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
			newTable(
				newTableGrid(),
				newTableRow(headerCell("a1"), headerCell("b1")),
				r("a2", "b2"),
				r("a3", "b3"),
			),
		),
	)

	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestTableHeaderAnyLocation verifies that th elements apply header styling
// regardless of their position in the table (not only in the first row).
func TestTableHeaderAnyLocation(t *testing.T) {
	a := assert.New(t)

	// First cell of each row is a th (row header), first row has no th at all.
	html := `
	<table>
		<tr>
			<td>name</td>
			<td>value</td>
		</tr>
		<tr>
			<th>row1</th>
			<td>val1</td>
		</tr>
		<tr>
			<th>row2</th>
			<td>val2</td>
		</tr>
	</table>
  `

	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	cell := func(txt string) *etree.Element {
		return newTableCell(newPara(newRun(newText(txt))))
	}

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newTable(
				newTableGrid(),
				newTableRow(cell("name"), cell("value")),
				newTableRow(headerCell("row1"), cell("val1")),
				newTableRow(headerCell("row2"), cell("val2")),
			),
		),
	)

	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestTableNoHeader verifies that when the first row uses td (not th), no
// header style is applied.
func TestTableNoHeader(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<td>a1</td>
			<td>b1</td>
		</tr>
		<tr>
			<td>a2</td>
			<td>b2</td>
		</tr>
	</table>
  `

	doc := newTestDocument()
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
			newTable(
				newTableGrid(),
				r("a1", "b1"),
				r("a2", "b2"),
			),
		),
	)

	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestTableSubheader verifies that class="subheader" on th or td applies cyan
// background with bold white text.
func TestTableSubheader(t *testing.T) {
	a := assert.New(t)

	html := `
	<table>
		<tr>
			<th class="subheader">sub1</th>
			<th class="subheader">sub2</th>
		</tr>
		<tr>
			<th>hdr1</th>
			<td class="subheader">sub3</td>
		</tr>
	</table>
  `

	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		tag("w:body",
			newTable(
				newTableGrid(),
				newTableRow(subheaderCell("sub1"), subheaderCell("sub2")),
				newTableRow(headerCell("hdr1"), subheaderCell("sub3")),
			),
		),
	)

	a.Equal(expected.StripBody().ToString(), doc.document.StripBody().ToString())
}

// TestColCount verifies the row col counter is correct
func TestColCount(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	c := newCursor(doc, "main")
	c.startTable().
		startTableRow().
		startTableCell(0).endTableCell().
		startTableCell(0).endTableCell().
		startTableCell(0).endTableCell()
	if a.IsType(c.e.token, &tableRow{}) {
		a.Equal(3, c.e.token.(*tableRow).col)
	}
}
