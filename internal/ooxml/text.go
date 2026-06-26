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

	"github.com/beevik/etree"
)

func newSpace() *etree.Element {
	return tag("w:t", attr("xml:space", "preserve"), " ")
}

func newText(txt string) *etree.Element {
	return tag("w:t", strings.Join(strings.Fields(txt), " "))
}

// txtState is a bitmap indicating the state of text across runs
type txtState struct {
	isStarted       bool
	isEmpty         bool
	isWhitespace    bool
	hasPfxSpace     bool
	hasSfxSpace     bool
	hasSpaceBetween bool
}

// nextText increments to the next block of text
func (s *txtState) nextText(txt string) *txtState {
	switch {
	case txt == "":
		s.isEmpty, s.isWhitespace = true, false
		s.hasPfxSpace, s.hasSfxSpace = false, false

	case strings.TrimSpace(txt) == "":
		s.isEmpty, s.isWhitespace = false, true
		s.hasPfxSpace, s.hasSfxSpace = true, true

	default:
		s.isEmpty, s.isWhitespace = false, false
		s.hasPfxSpace = strings.HasPrefix(txt, " ") || strings.HasPrefix(txt, "\n")
		s.hasSfxSpace = strings.HasSuffix(txt, " ") || strings.HasSuffix(txt, "\n")
	}
	return s
}

// nextSpace increments to the next gap between text
func (s *txtState) nextSpace() *txtState {
	s.hasSpaceBetween = false
	s.isStarted = s.isStarted || !(s.isEmpty || s.isWhitespace)
	return s
}

func (s *txtState) reset() *txtState {
	s.isStarted = false
	s.isEmpty = false
	s.isWhitespace = false
	s.hasPfxSpace, s.hasSfxSpace = false, false
	s.hasSpaceBetween = false
	return s
}

func (s *txtState) needsPfxSpace() bool {
	return s.isStarted && s.hasPfxSpace && !s.hasSpaceBetween
}

func (c *cursor) addText(txt string) *cursor {
	if c.txtState.isEmpty || (c.txtState.isWhitespace && !c.txtState.isStarted) {
		return c
	}
	switch c.e.token.(type) {
	case *body, *tableCell:
		c.startPara().startRun().addText(txt)
	case *para, *listItem, *header, *tableCaption: // Para types
		c.startRun().addText(txt)
	case *link:
		c.startRun(runLink{}).addText(txt)
	case *run:
		// Prefix space
		if c.txtState.needsPfxSpace() {
			c.AddXML(newSpace())
			c.txtState.hasSpaceBetween = true
		}

		if c.txtState.isWhitespace {
			return c
		}

		// Add text
		c.AddXML(newText(txt))
		c.txtState.nextSpace()

		// Suffix space
		if c.txtState.hasSfxSpace {
			c.AddXML(newSpace())
			c.txtState.hasSpaceBetween = true
		}
	}
	return c
}

func (c *cursor) nextText(txt string) *cursor {
	c.txtState.nextText(txt)
	return c
}

func (c *cursor) resetText() *cursor {
	c.txtState.reset()
	return c
}
