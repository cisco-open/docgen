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
	"strings"

	"github.com/beevik/etree"
)

type link struct{}

func (*link) tokenType() {}

func newAnchor(target string, children ...etree.Token) *etree.Element {
	return tag("w:hyperlink", attr("w:anchor", target), children)
}

func newLink(id int, children ...etree.Token) *etree.Element {
	return tag("w:hyperlink", attr("r:id", fmt.Sprintf("rId%d", id)), children)
}

func newLinkEl(doc *Document, target string) *Element {
	if strings.Contains(target, ".") {
		id := doc.newHyperlinkRel(target)
		return &Element{e: newLink(id), token: &link{}}
	}
	return &Element{e: newAnchor(target), token: &link{}}
}

func (c *cursor) startLink(target string) *cursor {
	switch c.e.token.(type) {
	case *body, *tableCell:
		c.startPara().AddChild(newLinkEl(c.doc, target))
	case *para, *header, *listItem:
		c.AddChild(newLinkEl(c.doc, target))
	case *run:
		c.endRun().AddChild(newLinkEl(c.doc, target))
	}
	return c
}

func (c *cursor) endLink() *cursor {
	return endToken[*link](c)
}
