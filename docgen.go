// Package docgen provides functionality for generating Microsoft Word documents
// from templates with HTML content and template tags.
//
// The package enables programmatic manipulation of Word documents in the
// Office Open XML (OOXML) format. It supports inserting HTML content,
// rendering template tags, and writing the modified document back to disk.
//
// Example usage:
//
//	template, err := os.ReadFile("template.docx")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	doc, err := docgen.NewDocument(template)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Insert HTML at a bookmark
//	htmlContent := strings.NewReader("<p>Hello, World!</p>")
//	err = doc.InsertHTML(htmlContent, "main")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Write the document
//	err = doc.WriteFile("output.docx")
//	if err != nil {
//	    log.Fatal(err)
//	}

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
package docgen

import (
	"io"

	"github.com/ciscotools/docgen/internal/ooxml"
)

// Document represents a Microsoft Word document in Office Open XML format.
// It provides methods for manipulating the document content, including
// inserting HTML, rendering template tags, and writing the document to disk.
//
// A Document is created using NewDocument with a template document in .docx format.
// Once created, you can modify the document using its methods and then write
// it back to disk using WriteFile or Write.
type Document struct {
	doc *ooxml.Document
}

// NewDocument creates a new Word document from a template.
//
// The template parameter should contain the bytes of a valid .docx file
// (Microsoft Word document in Office Open XML format). The template file
// serves as the base document which can then be modified using the returned
// Document's methods.
//
// The template must contain valid OOXML structure including:
//   - word/document.xml - the main document content
//   - word/_rels/document.xml.rels - relationship definitions
//   - word/numbering.xml - numbering definitions for lists
//
// Optional components that will be loaded if present:
//   - word/header1.xml - document header
//   - word/footer1.xml - document footer
//
// Returns a pointer to a Document that can be used to manipulate the
// Word document, or an error if the template is invalid or cannot be parsed.
//
// Example:
//
//	template, err := os.ReadFile("template.docx")
//	if err != nil {
//	    return err
//	}
//
//	doc, err := docgen.NewDocument(template)
//	if err != nil {
//	    return err
//	}
//
//	// Use the document...
//	err = doc.WriteFile("output.docx")
//	return err
func NewDocument(template []byte) (Document, error) {
	doc, err := ooxml.NewDocument(template)
	return Document{doc}, err
}

// WriteFile writes the document to a file at the specified path.
//
// This method serializes the entire Word document (including all modifications)
// back to the Office Open XML format and writes it to disk. The resulting file
// will be a valid .docx file that can be opened in Microsoft Word or other
// compatible applications.
//
// The path parameter specifies where to write the file. If the file already
// exists, it will be overwritten. Parent directories must exist.
//
// Returns an error if the file cannot be created or if there are issues
// serializing the document.
//
// Example:
//
//	err := doc.WriteFile("output.docx")
//	if err != nil {
//	    log.Fatal("Failed to write document:", err)
//	}
func (doc Document) WriteFile(path string) error {
	return doc.doc.WriteFile(path)
}

// Write writes the document to the provided io.Writer.
//
// This method serializes the entire Word document (including all modifications)
// to the Office Open XML format and writes it to the provided writer. This is
// useful for streaming the document to HTTP responses, buffers, or other
// destinations without writing to disk.
//
// The writer will receive a complete .docx file (ZIP archive containing XML files).
//
// Returns an error if there are issues serializing or writing the document.
//
// Example:
//
//	var buf bytes.Buffer
//	err := doc.Write(&buf)
//	if err != nil {
//	    log.Fatal("Failed to write document:", err)
//	}
//	// buf now contains the complete .docx file
func (doc *Document) Write(w io.Writer) error {
	return doc.doc.Write(w)
}

// InsertHTML inserts HTML content into the document at the specified bookmark location.
//
// This method parses HTML from the provided reader and converts it to
// corresponding Word document elements (paragraphs, runs, formatting, etc.),
// inserting them at the bookmark location.
//
// The input parameter is a reader containing HTML content. The HTML will be
// parsed and converted to OOXML elements.
//
// The bookmark parameter specifies where in the document to insert the HTML.
// The bookmark must exist in the template document (created using Insert > Bookmark
// in Microsoft Word). If the bookmark is not found, an error is returned.
//
// Supported HTML elements include:
//   - Text formatting: <b>, <i>, <u>, <s>, <sub>, <sup>
//   - Paragraphs: <p>, <div>
//   - Lists: <ul>, <ol>, <li>
//   - Tables: <table>, <tr>, <td>, <th>
//   - Links: <a href="...">
//   - Images: <img src="...">
//   - Headings: <h1> through <h6>
//   - Line breaks: <br>
//
// Returns an error if the bookmark is not found or if there are issues
// parsing the HTML or modifying the document.
//
// Example:
//
//	html := strings.NewReader("<p>This is <b>bold</b> text.</p>")
//	err := doc.InsertHTML(html, "content")
//	if err != nil {
//	    log.Fatal("Failed to insert HTML:", err)
//	}
func (doc Document) InsertHTML(input io.Reader, bookmark string) error {
	return doc.doc.InsertHTML(input, bookmark)
}

// InsertMarkdown inserts Markdown content into the document at the specified bookmark location.
//
// This method parses Markdown from the provided reader, converts it to HTML using
// goldmark, and then converts the HTML to corresponding Word document elements
// (paragraphs, runs, formatting, etc.), inserting them at the bookmark location.
//
// The input parameter is a reader containing Markdown content. The Markdown will be
// parsed and converted to OOXML elements.
//
// The bookmark parameter specifies where in the document to insert the Markdown.
// The bookmark must exist in the template document (created using Insert > Bookmark
// in Microsoft Word). If the bookmark is not found, an error is returned.
//
// Supported Markdown elements include:
//   - Text formatting: **bold**, *italic*, ~~strikethrough~~
//   - Paragraphs and line breaks
//   - Lists: ordered and unordered
//   - Tables
//   - Links: [text](url)
//   - Images: ![alt](src)
//   - Headings: # through ######
//   - Code blocks and inline code
//
// Returns an error if the bookmark is not found or if there are issues
// parsing the Markdown or modifying the document.
//
// Example:
//
//	markdown := strings.NewReader("# Title\n\nThis is **bold** text.")
//	err := doc.InsertMarkdown(markdown, "content")
//	if err != nil {
//	    log.Fatal("Failed to insert Markdown:", err)
//	}
func (doc Document) InsertMarkdown(input io.Reader, bookmark string) error {
	return doc.doc.InsertMarkdown(input, bookmark)
}

// RenderTags renders all template tags in the document using the provided context.
//
// Template tags are text placeholders in the format {{tagName}} that appear
// anywhere in the document (body, headers, footers). This method finds all
// such tags and replaces them with values from the context.
//
// The ctx parameter is the template context, typically a map or struct
// containing the values to substitute. The context is evaluated using
// Go template syntax, so you can use any valid Go template expression
// within the {{ }} delimiters.
//
// Tags can appear in:
//   - Document body
//   - Headers
//   - Footers
//
// Tags are evaluated independently, so each tag is treated as a separate
// template. Complex expressions like {{.user.name}} or {{add .a .b}} are
// supported if you register the appropriate template functions.
//
// Returns an error if there are issues parsing or rendering any template tag.
//
// Example:
//
//	ctx := map[string]interface{}{
//	    "title": "My Document",
//	    "author": "John Doe",
//	    "date": "2024-01-01",
//	}
//	err := doc.RenderTags(ctx)
//	if err != nil {
//	    log.Fatal("Failed to render tags:", err)
//	}
func (doc Document) RenderTags(ctx any) error {
	return doc.doc.RenderTags(ctx)
}

// GetBookmark finds and returns the bookmark element with the specified name.
//
// Bookmarks are named locations in a Word document that can be used as
// insertion points or references. This method searches the document for
// a bookmark with the given name and returns an Element representing
// the location.
//
// The name parameter specifies the bookmark name to search for. Bookmark
// names are case-sensitive and must match exactly.
//
// Returns a pointer to an Element representing the bookmark location,
// or nil if the bookmark is not found or is not within the document body.
//
// Example:
//
//	bookmark := doc.GetBookmark("main")
//	if bookmark == nil {
//	    log.Fatal("Bookmark 'main' not found")
//	}
func (doc Document) GetBookmark(name string) *ooxml.Element {
	return doc.doc.GetBookmark(name)
}
