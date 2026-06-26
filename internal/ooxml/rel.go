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
	"path/filepath"
	"strconv"
	"strings"
)

func (doc *Document) nextRelID() int {
	doc.lastRelIDOnce.Do(func() {
		rels := doc.rels.doc.FindElements("/" + "/Relationship")
		for _, rel := range rels {
			id, err := strconv.Atoi(
				strings.TrimPrefix(rel.SelectAttrValue("Id", "rId0"), "rId"))
			if err != nil {
				continue
			}
			if id > doc.lastRelID {
				doc.lastRelID = id
			}
		}
	})
	doc.lastRelID++
	return doc.lastRelID
}

func (doc *Document) newHyperlinkRel(target string) int {
	id := doc.nextRelID()
	doc.rels.doc.SelectElement("Relationships").AddChild(
		tag("Relationship",
			attr("Id", fmt.Sprintf("rId%d", id)),
			attr("Type",
				"http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink",
			),
			attr("Target", target),
			attr("TargetMode", "External"),
		),
	)
	return id
}

func (doc *Document) newImageRel(path string) int {
	id := doc.nextRelID()
	// internal path for rel is media/image{relId}.{ext}
	target := fmt.Sprintf("media/image%d%s", id, filepath.Ext(path))
	doc.rels.doc.SelectElement("Relationships").AddChild(
		tag("Relationship",
			attr("Id", fmt.Sprintf("rId%d", id)),
			attr("Type",
				"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
			),
			attr("Target", target),
		),
	)
	return id
}
