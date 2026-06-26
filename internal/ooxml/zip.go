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

	"github.com/beevik/etree"
)

// readXMLZip reads etree.Document from a zip.Reader
func readXMLZip(zr *zip.Reader, path string) (ooxmlDoc, error) {
	doc := etree.NewDocument()
	for _, f := range zr.File {
		if f.Name == path {
			rc, err := f.Open()
			if err != nil {
				return ooxmlDoc{}, err
			}
			_, err = doc.ReadFrom(rc)
			if err != nil {
				return ooxmlDoc{}, err
			}
			break
		}
	}
	return ooxmlDoc{
		doc:  doc,
		path: path,
	}, nil
}

// writeXMLZip writes a etree.Document to a zip.Writer
func (doc *ooxmlDoc) writeXMLZip(zw *zip.Writer) error {
	w, err := zw.Create(doc.path)
	if err != nil {
		return err
	}
	_, err = doc.doc.WriteTo(w)
	return err
}
