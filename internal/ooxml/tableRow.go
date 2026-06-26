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
	"github.com/beevik/etree"
)

// col and vMerge are used for tracking rowspan
//
// The exclusion of rowspan rows in HTML means we need to track and render this cell anyway.
// We do this by keeping vMerges in the table, indexed by column (col).
// For each new row, we read the vMerge object to determine if any cells are merged
// and need to be rendered automatically.
type tableRow struct {
	col     int
	vMerges map[int]*vMerge
}

func (*tableRow) tokenType() {}

func newTableRow(children ...any) *etree.Element {
	return tag("w:tr", children...)
}

func newTableRowEl(vMerges map[int]*vMerge) *Element {
	return &Element{e: newTableRow(), token: &tableRow{
		vMerges: vMerges,
	}}
}

func (c *cursor) startTableRow() *cursor {
	switch t := c.e.token.(type) {
	case *table:
		c.AddChild(newTableRowEl(t.vMerges)).
			addTableCellvMerge()
	}
	return c
}

func (c *cursor) addTableCellvMerge() *cursor {
	tr, ok := c.e.token.(*tableRow)
	if !ok {
		return c
	}
	for {
		// Check if a rowspan exists and is still within range
		if vMerge, ok := tr.vMerges[tr.col]; ok && vMerge.of > 1 && vMerge.row < vMerge.of {
			// If so, render an empty span cell
			c.startTableCell(0, append(vMerge.props, tableCellRowspan{})...).
				endTableCell()
			vMerge.row++
			// Continue to next iteration
			continue
		}
		return c
	}
}

func (c *cursor) endTableRow() *cursor {
	return endToken[*tableRow](c)
}
