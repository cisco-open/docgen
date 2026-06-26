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

func TestHeader(t *testing.T) {
	a := assert.New(t)

	html := `<h1>Header 1</h1>`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		newParaEl(Props(paraStyle{"Heading1"})).AddXML(
			newRun(newText("Header 1")),
		).e,
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestHeaderNoNumber(t *testing.T) {
	a := assert.New(t)

	html := `<h1 class="nonumber">Header 1</h1>`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		newParaEl(Props(paraStyle{"Heading1-NoNumbers"})).AddXML(
			newRun(newText("Header 1")),
		).e,
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestHeaderAutoNumber(t *testing.T) {
	a := assert.New(t)

	html := `
		<h2>A</h2>
		<h1 class="autonumber">B</h1>
	`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		tag("w:body",
			tag("w:bookmarkStart", attr("name", "main")),
			newHeaderEl(2, Props()).AddXML(newRun(newText("A"))).e,
			newHeaderEl(2, Props()).AddXML(newRun(newText("B"))).e,
		),
	)
	a.Equal(expected.ToString(), doc.document.ToString())
}

func TestHeaderAutoNumberUp(t *testing.T) {
	a := assert.New(t)

	html := `
		<h2>A</h2>
		<h1 class="autonumber-up">B</h1>
	`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		tag("w:body",
			tag("w:bookmarkStart", attr("name", "main")),
			newHeaderEl(2, Props()).AddXML(newRun(newText("A"))).e,
			newHeaderEl(3, Props()).AddXML(newRun(newText("B"))).e,
		),
	)
	a.Equal(expected.ToString(), doc.document.ToString())
}

func TestHeaderAutoNumberDown(t *testing.T) {
	a := assert.New(t)

	html := `
		<h2>A</h2>
		<h1 class="autonumber-down">B</h1>
	`

	// Write to test doc
	doc := newTestDocument()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	expected := newOOXMLDocFromTag(
		tag("w:body",
			tag("w:bookmarkStart", attr("name", "main")),
			newHeaderEl(2, Props()).AddXML(newRun(newText("A"))).e,
			newHeaderEl(1, Props()).AddXML(newRun(newText("B"))).e,
		),
	)
	a.Equal(expected.ToString(), doc.document.ToString())
}

func TestInvalidHeader(t *testing.T) {
	a := assert.New(t)

	// Test to confirm fix for panic on invalid HTML
	html := `
	<h1>Header 1
	<h2>Header 1.1</h2>
	`

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
}
