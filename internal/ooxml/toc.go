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

import "github.com/beevik/etree"

func newToC() *etree.Element {
	return tag("w:sdt",
		tag("w:sdtContent",
			attr("xmlns:w", "http://schemas.openxmlformats.org/wordprocessingml/2006/main"),
			tag("w:p",
				tag("w:r",
					tag("w:fldChar", attr("w:fldCharType", "begin"), attr("w:dirty", "true")),
				),
				tag("w:r",
					tag("w:instrText", attr("xml:space", "preserve"), `TOC \o "1-3" \h \z \u`),
				),
				tag("w:r",
					tag("w:fldChar", attr("w:fldCharType", "separate")),
				),
				tag("w:r",
					tag("w:fldChar", attr("w:fldCharType", "end")),
				),
			),
		),
	)
}

func (c *cursor) addToC() *cursor {
	switch c.e.token.(type) {
	case *body:
		c.AddXML(newToC())
	}
	return c
}
