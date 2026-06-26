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
)

type header struct {
	props *PropSet
	lvl   int
}

func (*header) tokenType() {}

func newHeaderEl(lvl int, props *PropSet) *Element {
	if props.Empty() {
		props.Add(paraStyle{fmt.Sprintf("Heading%d", lvl)})
	}
	return wrapElement(newParaEl(props), &header{props, lvl})
}

func (c *cursor) startHeader(lvl int, props ...prop) *cursor {
	c.headerLevel = lvl
	switch c.e.token.(type) {
	case *body:
		c.AddChild(newHeaderEl(lvl, Props(props...)))
	case *para:
		c.endPara().startHeader(lvl, props...)
	case *header:
		c.endHeader().startHeader(lvl, props...)
	case *run:
		c.endRun().startHeader(lvl, props...)
	}
	return c
}

func (c *cursor) endHeader() *cursor {
	c = endToken[*header](c)
	if e := findParent[*header](c.e); e != nil && e.Parent() != nil {
		c.headerLevel = e.token.(*header).lvl
		c.e = e.Parent()
	}
	return c
}
