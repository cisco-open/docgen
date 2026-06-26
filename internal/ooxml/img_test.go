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

func TestDecodeImage(t *testing.T) {
	a := assert.New(t)

	doc := newTestDocument()
	img, err := newImg(doc, "testdata/base16.png", 0, 0)
	a.Nil(err)
	a.Equal(100, img.x)
	a.Equal(100, img.y)
	a.Equal("base16", img.name)
	a.Equal(1, img.relID)

	img, err = newImg(doc, "testdata/cisco.jpg", 0, 0)
	a.Nil(err)
	a.Equal(180, img.x)
	a.Equal(180, img.y)
	a.Equal("cisco", img.name)
	a.Equal(2, img.relID)
}

func TestImage(t *testing.T) {
	a := assert.New(t)

	html := `
		<p>before</p>
		<p>
			<img src="/testdata/cisco.jpg" />
		</p>
		<p>after</p>
		<p>
			before<img src="/testdata/cisco.jpg" width="20" height="20"/>after
		</p>
	`

	doc, _ := getTestDoc()
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
	a.NoError(writeTestDoc(doc, "image.docx"))

	// Write to test doc
	doc = newTestDocument()
	html = `<img src="/testdata/cisco.jpg" />`
	a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))

	img, err := newImg(doc, "/testdata/cisco.jpg", 0, 0)
	if a.Nil(err) {
		expected := newOOXMLDocFromTag(
			newPara(
				newRun(
					newImage(1, img),
				),
			),
		)

		a.Equal(expected.ToString(), doc.document.StripBody().ToString())
	}
}
