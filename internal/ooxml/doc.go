// Package ooxml provides OpenXML abstractions for editing MS Word documents.
// This package implements low-level operations on Office Open XML (OOXML) format
// documents, particularly Microsoft Word .docx files.

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
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/beevik/etree"
)

type ooxmlDoc struct {
	path string
	doc  *etree.Document
}

// Dump pretty prints the document to dst
func (doc ooxmlDoc) Dump(dst io.Writer) error {
	doc.doc.Indent(2)
	_, err := doc.doc.WriteTo(dst)
	return err
}

// StripBody strips the body tag for testing
func (doc ooxmlDoc) StripBody() ooxmlDoc {
	body := doc.doc.SelectElement("w:body")
	if body == nil {
		return doc
	}
	for _, child := range body.FindElements("./*") {
		if child.Tag == "bookmarkStart" || child.Tag == "bookmarkEnd" {
			continue
		}
		doc.doc.AddChild(child)
	}
	doc.doc.RemoveChildAt(body.Index())
	return doc
}

// ToString returns dump output as a string
func (doc ooxmlDoc) ToString() string {
	var b strings.Builder
	doc.doc.Indent(2)
	_, _ = doc.doc.WriteTo(&b) // strings.Builder.Write never returns an error
	return b.String()
}

// Document represents a Microsoft Word document in the Office Open XML format.
// It maintains references to the main document XML, relationships, numbering
// definitions, headers, footers, images, and template tags.
//
// For more details on the OOXML format, see: http://officeopenxml.com/anatomyofOOXML.php
type Document struct {
	document  ooxmlDoc
	rels      ooxmlDoc
	numbering ooxmlDoc
	headers   []ooxmlDoc
	footers   []ooxmlDoc

	images   map[string]*img
	template []byte
	tags     []Tag

	oListID                 string
	uListID                 string
	addNumberingClassesOnce sync.Once

	lastImageID     int
	lastImageIDOnce sync.Once

	lastBookmarkID     int
	lastBookmarkIDOnce sync.Once

	lastRelID     int
	lastRelIDOnce sync.Once

	lastNumID     int
	lastNumIndex  int
	lastNumIDOnce sync.Once
}

// NewDocument creates a new Word document from a template.
//
// The template parameter should contain the bytes of a valid .docx file.
// This function parses the OOXML structure, extracting the document content,
// relationships, numbering definitions, headers, and footers.
//
// Returns a Document that can be modified using its methods, or an error
// if the template is invalid or required components are missing.
func NewDocument(template []byte) (*Document, error) {
	// Read zip file
	src := bytes.NewReader(template)
	zr, err := zip.NewReader(src, src.Size())
	if err != nil {
		return nil, err
	}

	// Read key document components
	document, err := readXMLZip(zr, "word/document.xml")
	if err != nil {
		return nil, err
	}
	rels, err := readXMLZip(zr, "word/_rels/document.xml.rels")
	if err != nil {
		return nil, err
	}
	numbering, err := readXMLZip(zr, "word/numbering.xml")
	if err != nil {
		return nil, err
	}
	headers := []ooxmlDoc{}
	if header, err := readXMLZip(zr, "word/header1.xml"); err == nil {
		headers = append(headers, header)
	}
	footers := []ooxmlDoc{}
	if footer, err := readXMLZip(zr, "word/footer1.xml"); err == nil {
		footers = append(footers, footer)
	}

	doc := &Document{
		document:  document,
		rels:      rels,
		numbering: numbering,
		headers:   headers,
		footers:   footers,
		images:    map[string]*img{},
		template:  template,
		tags:      []Tag{},
		oListID:   "1",
		uListID:   "2",
	}
	doc.addNumberingClasses()
	doc.findTags(doc.document)
	for _, header := range headers {
		doc.findTags(header)
	}
	for _, footer := range footers {
		doc.findTags(footer)
	}
	return doc, nil
}

// WriteFile writes the document to a file at the specified path.
// The output is a complete .docx file that can be opened in Microsoft Word.
//
// Returns an error if the file cannot be created or written.
func (doc *Document) WriteFile(path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	return doc.Write(dst)
}

// Write writes the document to the provided io.Writer.
// The output is a complete .docx file in OOXML format.
//
// Returns an error if there are issues serializing or writing the document.
func (doc *Document) Write(w io.Writer) error {
	src := bytes.NewReader(doc.template)
	zr, err := zip.NewReader(src, src.Size())
	if err != nil {
		return err
	}
	zw := zip.NewWriter(w)
	defer zw.Close()

	// Create path-based index of docs
	ooxmlDocsByPath := map[string]ooxmlDoc{
		doc.document.path:  doc.document,
		doc.rels.path:      doc.rels,
		doc.numbering.path: doc.numbering,
	}
	for _, header := range doc.headers {
		ooxmlDocsByPath[header.path] = header
	}
	for _, footer := range doc.footers {
		ooxmlDocsByPath[footer.path] = footer
	}

	for _, f := range zr.File {
		// If file is in the index, write the modified version
		if ooxmlDoc, ok := ooxmlDocsByPath[f.Name]; ok {
			if err := ooxmlDoc.writeXMLZip(zw); err != nil {
				return err
			}
			continue
		}
		// Otherwise write the pre-existing version
		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()
		w, err := zw.Create(f.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	}

	for _, img := range doc.images {
		defer img.r.Close()
		fn := fmt.Sprintf("image%d%s", img.relID, filepath.Ext(img.path))
		w, err := zw.Create("word/media/" + fn)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, img.r); err != nil {
			return err
		}
	}
	return nil
}
