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

// listItem is a list listItem paragraph
type listItem struct{}

func (*listItem) tokenType() {}

func newListItem(id string, lvl int) *etree.Element {
	return newPara(
		tag("w:pPr",
			tag("w:pStyle", attr("w:val", "ListParagraph")),
			tag("w:numPr",
				tag("w:ilvl", attr("w:val", strconv.Itoa(lvl))),
				tag("w:numId", attr("w:val", id)),
			),
		),
	)
}

func newListItemEl(id string, lvl int) *Element {
	return &Element{e: newListItem(id, lvl), token: &listItem{}}
}

func (c *cursor) startListItem() *cursor {
	switch t := c.e.token.(type) {
	case *list:
		c.AddChild(newListItemEl(t.id, t.depth))
	}
	return c
}

func (c *cursor) endListItem() *cursor {
	return endToken[*listItem](c)
}
