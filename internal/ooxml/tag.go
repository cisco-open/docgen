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
	"github.com/ciscotools/docgen/internal/html"
)

// Tag is a text templating tag
//
// Example: {{tagName}}
//
// Tags are used for simple text replacement in parts of the document
// where HTML templating is not required or lacks features, e.g.
// headers, footers, complex cover pages.
type Tag struct {
	name string
	e1   *etree.Element
	i1   int
	e2   *etree.Element
	i2   int
}

// render replaces the tag text and surrounding brackets with new text.
// If the tag delimiters have been overwritten by a prior render, prefix and
// suffix default to empty to avoid reading stale indices.
func (t Tag) render(txt string) {
	e1Text := t.e1.Text()
	e2Text := t.e2.Text()

	pfx := ""
	if t.i1 < len(e1Text) && e1Text[t.i1] == '{' {
		pfx = e1Text[:t.i1]
	}
	sfx := ""
	if t.i2 < len(e2Text) && e2Text[t.i2] == '}' {
		sfx = e2Text[t.i2+1:]
	}

	switch {
	case t.e1 == t.e2:
		t.e1.SetText(pfx + txt + sfx)

	case t.e1.Parent() == t.e2.Parent():
		t.e1.SetText(pfx + txt)
		t.e2.SetText(sfx)
		r := t.e1.Parent()
		for i := t.e1.Index() + 1; i < t.e2.Index(); i++ {
			if e, ok := r.Child[i].(*etree.Element); ok && e.Tag == "t" {
				e.SetText("")
			}
		}
	default:
		t.e1.SetText(pfx + txt)
		t.e2.SetText(sfx)

		r1 := t.e1.Parent()
		for i := t.e1.Index() + 1; i < len(r1.Child); i++ {
			if e, ok := r1.Child[i].(*etree.Element); ok && e.Tag == "t" {
				e.SetText("")
			}
		}

		r2 := t.e2.Parent()
		for i := 0; i < t.e2.Index(); i++ {
			if e, ok := r2.Child[i].(*etree.Element); ok && e.Tag == "t" {
				e.SetText("")
			}
		}

		p := r1.Parent()
		for i := r1.Index() + 1; i < r2.Index(); i++ {
			if r, ok := p.Child[i].(*etree.Element); ok && r.Tag == "r" {
				for _, t := range r.FindElements("./" + "/w:t") {
					t.SetText("")
				}
			}
		}
	}
}

// tagParseState tracks parsing state for identifying tags
type tagParseState struct {
	name       strings.Builder
	openingCnt byte
	closingCnt byte
	inTag      bool
	tag        Tag
}

// reset resets the tag parser, e.g. on new paragraphs or invalid tag text
func (s *tagParseState) reset() {
	s.openingCnt = 0
	s.closingCnt = 0
	s.inTag = false
	s.name.Reset()
	s.tag = Tag{}
}

// findTags parses text elements to locate valid tags
func (doc *Document) findTags(ooxmlDoc ooxmlDoc) {
	s := tagParseState{}
	for _, p := range ooxmlDoc.doc.FindElements("/" + "/w:p") {
		for _, t := range p.FindElements("./" + "/w:t") {
			for i, b := range []byte(t.Text()) {
				switch b {
				case '{':
					switch {
					case s.inTag:
						if s.name.Len() > 0 {
							s.reset()
							s.openingCnt++
						}
					case s.openingCnt == 0:
						s.openingCnt++
						s.tag.e1 = t
						s.tag.i1 = i
					default:
						s.inTag = true
					}
				case '}':
					switch {
					case !s.inTag:
						s.reset()
					case s.closingCnt == 0:
						s.closingCnt++
					case s.name.Len() == 0:
						s.reset()
					default:
						s.tag.e2 = t
						s.tag.i2 = i
						s.tag.name = s.name.String()
						doc.tags = append(doc.tags, s.tag)
						s.reset()
					}
				default:
					switch {
					case !s.inTag:
						if s.openingCnt > 0 {
							s.reset()
						}
					default:
						s.name.WriteByte(b)
					}
				}
			}
		}
		s.reset()
	}
}

// RenderTags renders all template tags in the document using the provided context.
//
// Template tags are text placeholders in the format {{tagName}} found throughout
// the document (body, headers, footers). Each tag is evaluated as a Go template
// expression using the provided context.
//
// Returns an error if there are issues parsing or rendering any template tag.
func (doc *Document) RenderTags(ctx any) error {
	type tagWithIndex struct {
		tag   Tag
		index int
	}

	// Group single-element tags by element; render multi-element tags immediately.
	// Tags within the same element are rendered in reverse order to preserve index validity.
	elementTags := make(map[*etree.Element][]tagWithIndex)
	var orderedElements []*etree.Element

	renderTag := func(tag Tag) error {
		htmlDoc := html.New("tag", ctx)
		if err := htmlDoc.RegisterTemplate("tag", strings.NewReader("{{"+tag.name+"}}")); err != nil {
			return err
		}
		var b strings.Builder
		if err := htmlDoc.RenderHTML(&b); err != nil {
			return err
		}
		tag.render(b.String())
		return nil
	}

	for i, tag := range doc.tags {
		elem := tag.e1
		if tag.e1 == tag.e2 {
			if _, exists := elementTags[elem]; !exists {
				orderedElements = append(orderedElements, elem)
			}
			elementTags[elem] = append(elementTags[elem], tagWithIndex{tag, i})
		} else {
			if err := renderTag(tag); err != nil {
				return err
			}
		}
	}

	for _, elem := range orderedElements {
		tags := elementTags[elem]
		for i := len(tags) - 1; i >= 0; i-- {
			if err := renderTag(tags[i].tag); err != nil {
				return err
			}
		}
	}

	return nil
}
