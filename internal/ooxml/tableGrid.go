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
	"fmt"

	"github.com/beevik/etree"
)

type tableGrid struct {
	colWidths []int
	cols      int
}

func (*tableGrid) tokenType() {}

func (tg *tableGrid) width() (w int) {
	for _, colWidth := range tg.colWidths {
		w += colWidth
	}
	return
}

func newTableGrid(children ...any) *etree.Element {
	return tag("w:tblGrid", children...)
}

func newTableGridCol(w int) *etree.Element {
	return tag("w:gridCol",
		attr("w:w", fmt.Sprintf("%d", 50*w)))
}

func newTableGridEl() *Element {
	return &Element{e: newTableGrid(), token: &tableGrid{}}
}

func (c *cursor) startTableGrid() *cursor {
	switch t := c.e.token.(type) {
	case *table:
		return c.SetElement(t.grid)
	}
	return c
}

func (c *cursor) endTableGrid() *cursor {
	if e := findParent[*tableGrid](c.e); e != nil && e.Parent() != nil {
		t := e.token.(*tableGrid)
		// If widths don't add up to 100% pad col 0
		w := t.width()
		if w != 100 && len(t.colWidths) > 0 {
			t.colWidths[0] += 100 - w
			e.e.RemoveChildAt(0)
			e.e.InsertChildAt(0, newTableGridCol(t.colWidths[0]))
		}
		c.e = e.Parent()
	}
	return c
}

func (c *cursor) addTableGridCol(w int) *cursor {
	switch t := c.e.token.(type) {
	case *tableGrid:
		t.colWidths = append(t.colWidths, w)
		t.cols++
		c.AddXML(newTableGridCol(w))
	}
	return c
}
