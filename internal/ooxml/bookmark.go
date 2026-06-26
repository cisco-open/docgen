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
	"strconv"

	"github.com/beevik/etree"
)

func (doc *Document) nextBookmarkID() int {
	doc.lastBookmarkIDOnce.Do(func() {
		bookmarks := doc.document.doc.FindElements("/" + "/w:bookmarkStart")
		for _, bookmark := range bookmarks {
			id := bookmark.SelectAttr("id")
			if id == nil {
				continue
			}
			i, err := strconv.Atoi(id.Value)
			if err != nil {
				continue
			}
			if i > doc.lastBookmarkID {
				doc.lastBookmarkID = i
			}
		}
	})
	doc.lastBookmarkID++
	return doc.lastBookmarkID
}

func newBookmarkStart(id int, name string) *etree.Element {
	return tag("w:bookmarkStart",
		attr("w:id", fmt.Sprintf("%d", id)),
		attr("w:name", name),
	)
}

func newBookmarkEnd(id int) *etree.Element {
	return tag("w:bookmarkEnd", attr("w:id", fmt.Sprintf("%d", id)))
}

func (c *cursor) addBookmark(name string) *cursor {
	switch c.e.token.(type) {
	case *body:
		id := c.doc.nextBookmarkID()
		c.AddXML(newBookmarkStart(id, name))
		c.AddXML(newBookmarkEnd(id))
	}
	return c
}

// GetBookmark finds and returns the bookmark element with the specified name.
//
// Returns a pointer to an Element representing the bookmark location, or nil
// if the bookmark is not found or is not within the document body.
func (doc *Document) GetBookmark(name string) *Element {
	query := fmt.Sprintf("/"+"/w:bookmarkStart[@name='%s']", name)
	e := doc.document.doc.FindElement(query)
	for ; e != nil; e = e.Parent() {
		parent := e.Parent()
		if parent == nil {
			return nil
		}
		if parent.Tag == "body" {
			return newBody(e)
		}
	}
	return nil
}

