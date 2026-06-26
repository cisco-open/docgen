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
	_ "embed"
	"strconv"

	"github.com/beevik/etree"
)

//go:embed xml/olistAbstractNum.xml
var olistAbstractNum []byte

//go:embed xml/ulistAbstractNum.xml
var ulistAbstractNum []byte

// list is an ordered or unordered list
type list struct {
	id    string
	depth int
}

func (*list) tokenType() {}

func (doc *Document) nextNumID() int {
	doc.lastNumIDOnce.Do(func() {
		nums := doc.numbering.doc.FindElements("/" + "/w:num")
		if len(nums) > 0 {
			doc.lastNumIndex = nums[len(nums)-1].Index()
		}
		for _, e := range nums {
			numID := e.SelectAttr("w:numId")
			if numID == nil {
				continue
			}
			id, err := strconv.Atoi(numID.Value)
			if err != nil {
				continue
			}
			if id > doc.lastNumID {
				doc.lastNumID = id
			}
		}
	})
	doc.lastNumID++
	return doc.lastNumID
}

// addNumberingClasses adds abstract number definitions
func (doc *Document) addNumberingClasses() {
	doc.addNumberingClassesOnce.Do(func() {
		addNumbering := func(root *etree.Element, i int, id int, tpl []byte) string {
			child := etree.NewDocument()
			child.ReadFromBytes(tpl)
			numEl := child.SelectElement("w:abstractNum")
			if numEl == nil {
				return strconv.Itoa(id)
			}
			numEl.CreateAttr("w:abstractNumId", strconv.Itoa(id))
			root.InsertChildAt(i, numEl)
			return strconv.Itoa(id)
		}
		abstractNums := doc.numbering.doc.FindElements("/" + "/w:abstractNum")
		if len(abstractNums) == 0 {
			return
		}
		lastID := -1
		for _, abstractNum := range abstractNums {
			idAttr := abstractNum.SelectAttr("w:abstractNumId")
			if idAttr == nil {
				continue
			}
			id, err := strconv.Atoi(idAttr.Value)
			if err != nil {
				continue
			}
			if id > lastID {
				lastID = id
			}
		}
		lastChild := abstractNums[len(abstractNums)-1]
		i := lastChild.Index()
		root := lastChild.Parent()
		doc.oListID = addNumbering(root, i+1, lastID+1, olistAbstractNum)
		doc.uListID = addNumbering(root, i+2, lastID+2, ulistAbstractNum)
	})
}

// newListEl wraps the current element turning it into a list
func newListEl(doc *Document, root *Element, kind string) *Element {
	// Check if this is a child of an existing list
	id := strconv.Itoa(doc.nextNumID())

	// https://stackoverflow.com/questions/58622437/purpose-of-abstractnum-and-numberinginstance
	if e := doc.numbering.doc.SelectElement("w:numbering"); e != nil {
		e.InsertChildAt(doc.lastNumIndex+1,
			tag("w:num", attr("w:numId", id),
				tag("w:abstractNumId", attr("w:val", kind)),
				tag("w:lvlOverride", attr("w:ilvl", "0"),
					tag("w:startOverride", attr("w:val", "1")),
				),
			),
		)
		doc.lastNumIndex++
	}
	return wrapElement(root, &list{id: id})
}

func (c *cursor) startList(kind string) *cursor {
	switch t := c.e.token.(type) {
	case *body, *tableCell:
		c.SetElement(newListEl(c.doc, c.e, kind))
	case *list:
		t.depth++
	case *para:
		c.SetElement(c.e.Parent()).startList(kind)
	case *run:
		c.SetElement(c.e.Parent().Parent()).startList(kind)
	}
	return c
}

func (c *cursor) endList() *cursor {
	if e := findParent[*list](c.e); e != nil {
		if lst := e.token.(*list); lst.depth > 0 {
			lst.depth--
			c.SetElement(e)
		} else {
			c.SetElement(c.e.Parent())
		}
	}
	return c
}
