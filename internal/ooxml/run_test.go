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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunPropsBasic(t *testing.T) {
	a := assert.New(t)

	html := `plain<b>bold</b>`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("plain")),
			newRun(tag("w:rPr", tag("w:b")), newText("bold")),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestRunPropsComprehensive(t *testing.T) {
	a := assert.New(t)

	html := `
    plain
    <b>bold</b>
    <i>italic</i>
    <b>bold<i>bold-italic</i>bold</b>
  `

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("plain"), newSpace()),
			newRun(tag("w:rPr", tag("w:b")), newText("bold")),
			newRun(newSpace()),
			newRun(tag("w:rPr", tag("w:i")), newText("italic")),
			newRun(newSpace()),
			newRun(tag("w:rPr", tag("w:b")), newText("bold")),
			newRun(tag("w:rPr", tag("w:i"), tag("w:b")), newText("bold-italic")),
			newRun(tag("w:rPr", tag("w:b")), newText("bold")),
			newRun(newSpace()),
		),
	)

	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}
