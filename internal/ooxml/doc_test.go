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
	"os"
	"path/filepath"
	"testing"

	"github.com/beevik/etree"
	"github.com/cisco-open/docgen/internal/templates"
	"github.com/stretchr/testify/assert"
)

const testOutputPath = "testoutput"

func init() {
	os.MkdirAll(testOutputPath, 0o700)
}

func newTestDocument() *Document {
	document := etree.NewDocument()
	document.ReadFromString(`
    <w:body>
      <w:bookmarkStart name="main" />
    </w:body>
  `)
	numbering := etree.NewDocument()
	numbering.ReadFromString(`
    <w:numbering>
    </w:numbering>
  `)
	rels := etree.NewDocument()
	rels.ReadFromString(`
		<Relationships>
		</Relationships>
	`)
	return &Document{
		document:  ooxmlDoc{doc: document},
		numbering: ooxmlDoc{doc: numbering},
		rels:      ooxmlDoc{doc: rels},
		images:    map[string]*img{},
	}
}

func (doc *Document) updateDocument(tag *etree.Element) {
	doc.document = newOOXMLDocFromTag(tag)
	doc.tags = []Tag{}
	doc.findTags(doc.document)
}

func newOOXMLDocFromTag(e *etree.Element) ooxmlDoc {
	doc := etree.NewDocument()
	doc.AddChild(e)
	return ooxmlDoc{doc: doc}
}

func getTestDoc() (*Document, error) {
	doc, err := NewDocument(templates.TestTemplate)
	if err != nil {
		return nil, err
	}
	doc.addNumberingClasses()
	return doc, nil
}

func writeTestDoc(doc *Document, name string) error {
	path := filepath.Join(testOutputPath, name)
	err := doc.WriteFile(path)
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	return err
}

func TestReadWriteDocument(t *testing.T) {
	a := assert.New(t)
	doc, err := getTestDoc()
	a.NoError(err)
	a.NoError(writeTestDoc(doc, "docrw.docx"))
}
