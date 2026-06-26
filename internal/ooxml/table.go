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

type vMerge struct {
	row   int
	of    int
	props []prop
}

type table struct {
	grid    *Element
	vMerges map[int]*vMerge
}

func (*table) tokenType() {}

func newTable(children ...*etree.Element) *etree.Element {
	return tag("w:tbl",
		tag("w:tblPr",
			tag("w:tblStyle", attr("w:val", "CiscoCXTableDefault")),
			tag("w:tblW", attr("w:w", "5000"), attr("w:type", "pct")),
			// Disable all conditional formatting from the table style so that
			// header/subheader styling is applied explicitly per-cell.
			tag("w:tblLook",
				attr("w:val", "0000"),
				attr("w:firstRow", "0"),
				attr("w:lastRow", "0"),
				attr("w:firstColumn", "0"),
				attr("w:lastColumn", "0"),
				attr("w:noHBand", "0"),
				attr("w:noVBand", "0"),
			),
		),
		children,
	)
}

func newTableEl() *Element {
	e := &Element{e: newTable(), token: &table{
		vMerges: map[int]*vMerge{},
	}}
	e.token.(*table).grid = e.AddChild(newTableGridEl())
	return e
}

func (c *cursor) startTable() *cursor {
	switch c.e.token.(type) {
	case *body:
		c.AddChild(newTableEl())
	default:
		if body := findParent[*body](c.e); body != nil {
			c.SetElement(body).AddChild(newTableEl())
		}
	}
	return c
}

func (c *cursor) endTable() *cursor {
	return endToken[*table](c)
}
