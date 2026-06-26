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

type run struct {
	props *PropSet
}

func (*run) tokenType() {}

func newRun(children ...any) *etree.Element {
	return tag("w:r", children...)
}

func newRunEl(props *PropSet) *Element {
	r := newElement(newRun(), &run{props})
	if props.Empty() {
		return r
	}
	// Apply row styles
	rPr := tag("w:rPr")
	for _, prop := range props.Iter() {
		prop.addProp(rPr)
	}
	r.e.AddChild(rPr)
	return r
}

type runBold struct{}

func (runBold) addProp(rPr *etree.Element) {
	rPr.AddChild(tag("w:b"))
}

type runItalic struct{}

func (runItalic) addProp(rPr *etree.Element) {
	rPr.AddChild(tag("w:i"))
}

type runColor struct {
	val        string
	themeColor string
}

func (r runColor) addProp(rPr *etree.Element) {
	e := tag("w:color", attr("w:val", r.val))
	if r.themeColor != "" {
		e.CreateAttr("w:themeColor", r.themeColor)
	}
	rPr.AddChild(e)
}

type runLink struct{}

func (runLink) addProp(rPr *etree.Element) {
	rPr.AddChild(tag("w:rStyle", attr("w:val", "Hyperlink")))
}

// cellRunProps returns any run-level props inherited from the nearest tableCell
// ancestor (e.g. bold + white text for header cells), or nil if none.
func cellRunProps(e *Element) []prop {
	if tc := findParent[*tableCell](e); tc != nil {
		return tc.token.(*tableCell).runProps
	}
	return nil
}

func (c *cursor) startRun(props ...prop) *cursor {
	// Merge any cell-inherited run props (e.g. bold+colour for th cells)
	startProps := Props(append(cellRunProps(c.e), props...)...)
	switch t := c.e.token.(type) {
	case *body, *tableCell:
		c.startPara().AddChild(newRunEl(startProps))
	case *para, *listItem, *header, *tableCaption: // Para types
		c.AddChild(newRunEl(startProps))
	case *link:
		c.AddChild(newRunEl(startProps.Add(runLink{})))
	case *run:
		if startProps.Equals(t.props) {
			return c
		}
		c.endRun(t.props.Iter()...).
			AddChild(newRunEl(startProps.Union(t.props)))
	}
	return c
}

func (c *cursor) endRun(props ...prop) *cursor {
	endProps := Props(props...)
	if e := findParent[*run](c.e); e != nil {
		c.e = e.Parent()
		thisProps := e.token.(*run).props
		nextProps := thisProps.Difference(endProps)
		// Start a new run with any remaining props
		if !nextProps.Empty() {
			return c.AddChild(newRunEl(nextProps))
		}
	}
	return c
}
