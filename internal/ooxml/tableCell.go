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

	"github.com/beevik/etree"
)

type tableCell struct {
	props    *PropSet
	runProps []prop // run-level props inherited by all runs in this cell
}

func (*tableCell) tokenType() {}

func newTableCell(children ...any) *etree.Element {
	return tag("w:tc", children...)
}

func newTableCellEl(props *PropSet) *Element {
	// Collect run-level props contributed by cell-level marker props.
	var runProps []prop
	for _, p := range props.Iter() {
		if cp, ok := p.(cellRunPropContributor); ok {
			runProps = append(runProps, cp.runProps()...)
		}
	}
	tc := &Element{e: newTableCell(), token: &tableCell{props, runProps}}
	if props.Empty() {
		return tc
	}
	tcPr := tag("w:tcPr")
	for _, prop := range props.Iter() {
		prop.addProp(tcPr)
	}
	// Only attach w:tcPr if at least one prop wrote a child element
	if len(tcPr.Child) > 0 {
		tc.e.AddChild(tcPr)
	}
	return tc
}

// cellRunPropContributor is implemented by cell props that also inject run-level
// formatting (e.g. bold + text colour) into every run created within the cell.
type cellRunPropContributor interface {
	runProps() []prop
}

// tableCellHeader marks a <th> cell.
// Cell background: navy (1E4471). Run style: bold + white text.
type tableCellHeader struct{}

func (tableCellHeader) addProp(tcPr *etree.Element) {
	tcPr.AddChild(tag("w:shd",
		attr("w:val", "clear"),
		attr("w:color", "auto"),
		attr("w:fill", "1E4471"),
		attr("w:themeFill", "text2"),
	))
}

func (tableCellHeader) runProps() []prop {
	return []prop{runBold{}, runColor{val: "FFFFFF", themeColor: "background2"}}
}

// tableCellSubheader marks a <th class="subheader"> or <td class="subheader"> cell.
// Cell background: cyan (00BCEB). Run style: bold + white text.
type tableCellSubheader struct{}

func (tableCellSubheader) addProp(tcPr *etree.Element) {
	tcPr.AddChild(tag("w:shd",
		attr("w:val", "clear"),
		attr("w:color", "auto"),
		attr("w:fill", "00BCEB"),
		attr("w:themeFill", "accent1"),
	))
}

func (tableCellSubheader) runProps() []prop {
	return []prop{runBold{}, runColor{val: "FFFFFF", themeColor: "background2"}}
}

type tableCellColspan struct{ cols int }

func (t tableCellColspan) addProp(tcPr *etree.Element) {
	tcPr.AddChild(tag("w:gridSpan", attr("w:val", strconv.Itoa(t.cols))))
}

type tableCellRowspan struct {
	start bool
}

func (tcr tableCellRowspan) addProp(tcPr *etree.Element) {
	if tcr.start {
		tcPr.AddChild(tag("w:vMerge", attr("w:val", "restart")))
	} else {
		tcPr.AddChild(tag("w:vMerge"))
	}
}

func (c *cursor) startTableCell(rowspan int, props ...prop) *cursor {
	startProps := Props(props...)
	switch t := c.e.token.(type) {
	case *tableRow:
		if rowspan > 1 {
			t.vMerges[t.col] = &vMerge{row: 1, of: rowspan, props: props}
			startProps.Add(tableCellRowspan{true})
		}
		t.col++
		c.AddChild(newTableCellEl(startProps))
	}
	return c
}

func (c *cursor) endTableCell() *cursor {
	if e := findParent[*tableCell](c.e); e != nil && e.Parent() != nil {
		// An empty w:tc must have at least one child element
		if e.e.SelectElement("w:p") == nil {
			e.AddXML(newPara())
		}
		c.e = e.Parent()
		// After ending the cell, we check vMerge again for the next cell
		c.addTableCellvMerge()
	}
	return c
}
