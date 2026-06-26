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

type cursor struct {
	doc         *Document
	e           *Element
	txtState    txtState
	headerLevel int
}

// newCursor creates a new cursor
func newCursor(doc *Document, bookmark string) *cursor {
	return &cursor{
		doc: doc,
		e:   doc.GetBookmark(bookmark),
	}
}

// AddChild adds a child element to this cursor
func (c *cursor) AddChild(e *Element) *cursor {
	return c.SetElement(c.e.AddChild(e))
}

// AddXML adds a child etree.Element to this cursor
func (c *cursor) AddXML(e *etree.Element) *cursor {
	c.e.AddXML(e)
	return c
}

// SetElement sets the current element
func (c *cursor) SetElement(e *Element) *cursor {
	c.e = e
	return c
}

// findParent finds the parent element of type T
func findParent[T token](e *Element) *Element {
	for {
		if e == nil {
			return nil
		}
		if _, ok := e.token.(T); ok {
			return e
		}
		e = e.Parent()
	}
}

// endToken finds the parent element and returns its parent
func endToken[T token](c *cursor) *cursor {
	if e := findParent[T](c.e); e != nil && e.Parent() != nil {
		c.e = e.Parent()
	}
	return c
}
