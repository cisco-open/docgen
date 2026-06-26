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

	"github.com/beevik/etree"
)

type para struct {
	props *PropSet
}

func (*para) tokenType() {}

func newPara(children ...any) *etree.Element {
	return tag("w:p", children...)
}

func newParaEl(props *PropSet) *Element {
	p := &Element{e: newPara(), token: &para{props}}
	if props.Empty() {
		return p
	}
	pPr := tag("w:pPr")
	for _, prop := range props.Iter() {
		prop.addProp(pPr)
	}
	p.e.AddChild(pPr)
	return p
}

type paraStyle struct{ style string }

func (prop paraStyle) addProp(pPr *etree.Element) {
	pPr.AddChild(tag("w:pStyle",
		attr("w:val", strings.ReplaceAll(prop.style, " ", "")),
	))
}

func (c *cursor) startPara(props ...prop) *cursor {
	switch c.e.token.(type) {
	case *body, *tableCell:
		c.AddChild(newParaEl(Props(props...)))
	case *para:
		c.endPara().startPara()
	}
	return c
}

func (c *cursor) endPara() *cursor {
	return endToken[*para](c)
}
