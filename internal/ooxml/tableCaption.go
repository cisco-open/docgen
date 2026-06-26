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

// tableCaption is a "Caption" styled paragraph
type tableCaption struct{}

func (tableCaption) tokenType() {}

func newTableCaptionEl() *Element {
	p := newParaEl(Props(paraStyle{"Caption"}))
	p.token = &tableCaption{}
	return p
}

// startTableCaption inserts a caption *BEFORE* the current table element
func (c *cursor) startTableCaption() *cursor {
	switch c.e.token.(type) {
	case *table:
		caption := newTableCaptionEl()
		tbl := c.e
		bdy := findParent[*body](tbl)

		// Set table as parent of caption
		// even though caption is technically rendered outside of table
		caption.parent = tbl

		// Insert element *before* the table
		bdy.e.InsertChildAt(tbl.e.Index(), caption.e)

		// Update the body index to point to the table's new index
		bdy.token.(*body).index = tbl.e.Index()

		// Set caption as current element
		c.SetElement(caption)
	}
	return c
}

func (c *cursor) endTableCaption() *cursor {
	return endToken[*tableCaption](c)
}
