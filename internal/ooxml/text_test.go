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

func TestTextSpacing(t *testing.T) {
	a := assert.New(t)
	compare := func(html string, vals ...string) {
		doc := newTestDocument()
		a.NoError(doc.InsertHTML(strings.NewReader(html), "main"))
		els := doc.document.doc.FindElements("/" + "/w:t")
		a.Len(els, len(vals))
		txts := []string{}
		for _, e := range els {
			txts = append(txts, e.Text())
		}
		a.Equal(vals, txts)
	}

	// Space outside tag
	compare(
		`one <b>two</b> three`,
		"one", " ", "two", " ", "three",
	)
	// Space inside tag
	compare(
		`one<b> two </b>three`,
		"one", " ", "two", " ", "three",
	)
	// Extra spaces
	compare(
		`one  <b> two  </b>  three</p>`,
		"one", " ", "two", " ", "three",
	)
	// Inside paragraph
	compare(
		`<p>one <b>two</b> three</p>`,
		"one", " ", "two", " ", "three",
	)
	// Multiple lines and spaces
	compare(
		` one
		    <b> two </b>
		  three
		`,
		"one", " ", "two", " ", "three",
	)
}

func TestTextStateEmpty(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("")
	a.True(s.isEmpty)
	a.False(s.isWhitespace)
}

func TestTextStateWhitespace(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText(" ")
	a.False(s.isEmpty)
	a.True(s.isWhitespace)
}

func TestTextStateText(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("text")
	a.False(s.isEmpty)
	a.False(s.isWhitespace)
}

func TestTextStateStartingSpace(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("   \n ")
	a.False(s.isStarted)
	a.True(s.isWhitespace)
}

func TestTextStateIsStarted(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("text")
	a.False(s.isStarted)
	s.nextSpace()
	a.True(s.isStarted)
}

func TestTextStateNoSpaces(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("text")
	a.False(s.hasPfxSpace)
	a.False(s.hasSfxSpace)
}

func TestTextStateMultipleTexts(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText(" text ").nextText("text")
	a.False(s.hasPfxSpace)
	a.False(s.hasSfxSpace)
}

func TestTextStateTrailingEmpty(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("text").nextText("")
	a.False(s.hasPfxSpace)
	a.False(s.hasSfxSpace)
}

func TestTextStatePfxSpace(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText(" text")
	a.True(s.hasPfxSpace)
}

func TestTextStateSfxSpace(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.nextText("text ")
	a.True(s.hasSfxSpace)
}

func TestTextStatePfxSfxSpaces(t *testing.T) {
	a := assert.New(t)
	s := txtState{}
	s.reset().nextText(" text ")
	a.True(s.hasPfxSpace)
	a.True(s.hasSfxSpace)

	s.reset().nextText(" ")
	a.True(s.hasPfxSpace)
	a.True(s.hasSfxSpace)

	s.reset().nextText("\n")
	a.True(s.hasPfxSpace)
	a.True(s.hasSfxSpace)
}
