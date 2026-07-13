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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cisco-open/docgen/internal/log"

	"golang.org/x/net/html"
)

const (
	defaultParaStyle = "Normal-6ptspacing"
)

// htmlToken is a parsed HTML tag with lazy attribute loading.
type htmlToken struct {
	tag      string
	hasAttrs bool
	attrs    map[string]string
	z        *html.Tokenizer
}

func newToken(z *html.Tokenizer) *htmlToken {
	name, hasAttrs := z.TagName()
	return &htmlToken{
		tag:      string(name),
		hasAttrs: hasAttrs,
		z:        z,
	}
}

func (t *htmlToken) attr(key string) string {
	if !t.hasAttrs {
		return ""
	}
	// Read attrs into map
	if t.attrs == nil {
		t.attrs = map[string]string{}
		for {
			k, v, more := t.z.TagAttr()
			if len(k) > 0 {
				t.attrs[string(k)] = string(v)
			}
			if !more {
				break
			}
		}
	}
	if v, ok := t.attrs[key]; ok {
		return v
	}
	return ""
}

func parseColWidth(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty style")
	}
	s = strings.TrimPrefix(s, "width: ")
	s = strings.TrimSuffix(s, "%")
	return strconv.Atoi(s)
}

// InsertHTML inserts HTML content into the document at the specified bookmark location.
//
// This method parses HTML from the input reader and converts it to corresponding
// OOXML elements, inserting them at the bookmark position. The bookmark must exist
// in the document.
//
// Returns an error if the bookmark is not found or if there are issues parsing the HTML.
func (doc *Document) InsertHTML(input io.Reader, bookmark string) error {
	c := newCursor(doc, bookmark)
	if c.e == nil {
		return fmt.Errorf("bookmark %s not found", bookmark)
	}
	z := html.NewTokenizer(input)
	for {
		switch tt := z.Next(); tt {

		case html.ErrorToken:
			return nil

		case html.TextToken:
			txt := string(z.Text())
			c.nextText(txt).addText(txt)

		case html.StartTagToken:
			switch tkn := newToken(z); tkn.tag {
			// Misc
			case "a":
				c.startLink(tkn.attr("href"))
			case "br":
				c.startPara().endPara()
			case "p", "div":
				// div is treated as an alternative <p> for backwards compatibility
				class := tkn.attr("class")
				switch class {
				case "pagebreak":
					c.addPageBreak()
				case "toc":
					c.addToC()
				default:
					style := tkn.attr("data-custom-style")
					props := []prop{}
					if style != "" {
						props = append(props, paraStyle{style})
					} else {
						props = append(props, paraStyle{defaultParaStyle})
					}
					c.resetText().startPara(props...)
				}

			// Styles
			case "h1", "h2", "h3", "h4", "h5", "h6":
				lvl, _ := strconv.Atoi(string(tkn.tag)[1:])
				id := tkn.attr("id")
				if id != "" {
					c.addBookmark("id")
				}
				props := []prop{}
				for cls := range strings.FieldsSeq(tkn.attr("class")) {
					switch cls {
					case "nonumber":
						props = append(props, paraStyle{fmt.Sprintf("Heading%d-NoNumbers", lvl)})
					case "autonumber":
					case "autonumber-up":
						c.headerLevel++
					case "autonumber-down":
						c.headerLevel--
					}
					if strings.HasPrefix(cls, "autonumber") {
						if c.headerLevel > 6 {
							c.headerLevel = 6
						}
						if c.headerLevel < 1 {
							c.headerLevel = 1
						}
						lvl = c.headerLevel
					}
				}
				c.resetText().startHeader(lvl, props...)
			case "b", "strong":
				c.startRun(runBold{})
			case "i", "em":
				c.startRun(runItalic{})

			// Lists
			case "ol":
				c.startList(c.doc.oListID)
			case "ul":
				c.startList(c.doc.uListID)
			case "li":
				c.resetText().startListItem()

			// Tables
			case "table":
				c.startTable()
			case "caption":
				c.startTableCaption()
			case "colgroup":
				c.startTableGrid()
			case "tr":
				c.startTableRow()
			case "th", "td":
				props := []prop{}
				if tkn.tag == "th" {
					if tkn.attr("class") == "subheader" {
						props = append(props, tableCellSubheader{})
					} else {
						props = append(props, tableCellHeader{})
					}
				} else if tkn.attr("class") == "subheader" {
					props = append(props, tableCellSubheader{})
				}
				if colspan := tkn.attr("colspan"); colspan != "" {
					if cols, err := strconv.Atoi(colspan); err == nil {
						props = append(props, tableCellColspan{cols})
					}
				}
				rowspan := 0
				if rs := tkn.attr("rowspan"); rs != "" {
					if rows, err := strconv.Atoi(rs); err == nil {
						rowspan = rows
					}
				}
				c.resetText().startTableCell(rowspan, props...)
			}

		case html.EndTagToken:
			switch tkn := newToken(z); tkn.tag {
			case "a":
				c.endLink()
			case "b":
				c.endRun(runBold{})
			case "i":
				c.endRun(runItalic{})
			case "h1", "h2", "h3", "h4", "h5", "h6":
				c.resetText().endHeader()
			case "li":
				c.resetText().endListItem()
			case "ol", "ul":
				c.endList()
			case "p", "div":
				c.resetText().endPara()

			// Table
			case "table":
				c.endTable()
			case "caption":
				c.endTableCaption()
			case "colgroup":
				c.endTableGrid()
			case "td", "th":
				c.resetText().endTableCell()
			case "tr":
				c.endTableRow()
			}

		case html.SelfClosingTagToken:
			switch tkn := newToken(z); tkn.tag {
			case "br":
				c.startPara().endPara()
			case "hr":
				// low-priority; not currently in use in any templates
			case "img":
				src := tkn.attr("src")
				x, _ := strconv.Atoi(tkn.attr("width"))
				y, _ := strconv.Atoi(tkn.attr("height"))
				c.addImage(src, x, y)
			case "div":
				switch tkn.attr("class") {
				case "pagebreak":
					c.addPageBreak()
				case "toc":
					c.addToC()
				}
			case "col":
				w, err := parseColWidth(tkn.attr("style"))
				if err != nil {
					log.Warn().Msgf("cannot parse col width for %s: %v", tkn.tag, err)
					continue
				}
				c.addTableGridCol(w)
			}
		}
	}
}
