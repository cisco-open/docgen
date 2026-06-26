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

// token is an abstract element type
type token interface {
	tokenType()
}

// Element is an OOXML Element
type Element struct {
	e      *etree.Element
	parent *Element
	token  token
}

// newElement creates a new element
func newElement(e *etree.Element, t token) *Element {
	return &Element{e: e, token: t}
}

// wrapElement wraps an existing element with a new token
func wrapElement(root *Element, t token) *Element {
	return &Element{
		e:      root.e,
		parent: root,
		token:  t,
	}
}

// Parent returns the Parent token or nil
func (e *Element) Parent() *Element {
	return e.parent
}

// AddChild adds a new child token and returns the child
func (e *Element) AddChild(child *Element) *Element {
	e.AddXML(child.e)
	child.parent = e
	return child
}

// AddXML adds a new child XML element
func (e *Element) AddXML(child *etree.Element) *Element {
	switch t := e.token.(type) {
	case *body:
		e.e.InsertChildAt(t.index+1, child)
		t.index = child.Index()
	case *list:
		if t, ok := e.Parent().token.(*body); ok {
			e.e.InsertChildAt(t.index+1, child)
			t.index = child.Index()
		} else {
			e.e.AddChild(child)
		}
	default:
		e.e.AddChild(child)
	}
	return e
}
