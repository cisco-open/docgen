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
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ciscotools/docgen/internal/log"

	"github.com/beevik/etree"
)

type img struct {
	path  string
	name  string
	relID int
	r     io.ReadSeekCloser
	x     int
	y     int
}

func (*img) tokenType() {}

func newImg(doc *Document, path string, x, y int) (*img, error) {
	// Remove pfx slash from relative HTML paths
	path = strings.TrimPrefix(path, "/")
	if img, ok := doc.images[path]; ok {
		return img, nil
	}

	// Normalize jpeg extension
	ooxmlPath := path
	if filepath.Ext(path) == ".jpg" {
		ooxmlPath = strings.TrimSuffix(ooxmlPath, ".jpg") + ".jpeg"
	}

	// Validate extension
	ext := filepath.Ext(ooxmlPath)
	if ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("unsupported image format for %s", path)
	}

	// Read image
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, _, _ := image.Decode(r)
	bounds := m.Bounds()
	if _, err := r.Seek(0, 0); err != nil {
		r.Close()
		return nil, err
	}
	if x == 0 {
		x = bounds.Max.X
	}
	if y == 0 {
		y = bounds.Max.Y
	}

	// create new rel ID
	relID := doc.newImageRel(ooxmlPath)

	img := &img{
		path:  ooxmlPath,
		name:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		relID: relID,
		r:     r,
		x:     x,
		y:     y,
	}
	doc.images[path] = img
	return img, nil
}

// newImage creates OOXML for an instance of an image
func newImage(id int, img *img) *etree.Element {
	relID := fmt.Sprintf("rId%d", img.relID)
	imgID := strconv.Itoa(id)

	// Dimensions are in EMU
	// 1in = 914400emu
	// Assuming 96 dpi, 1px = 9525emu
	cx := strconv.Itoa(img.x * 9525)
	cy := strconv.Itoa(img.y * 9525)
	return tag("w:drawing",
		tag("wp:inline",
			attr("distT", "0"),
			attr("distB", "0"),
			attr("distL", "0"),
			attr("distR", "0"),
			tag("wp:extent", attr("cx", cx), attr("cy", cy)),
			tag("wp:effectExtent",
				attr("l", "0"),
				attr("t", "0"),
				attr("r", "0"),
				attr("b", "0"),
			),
			tag("wp:docPr",
				attr("id", imgID),
				attr("name", img.name),
				attr("descr", img.name),
			),
			tag("wp:cNvGraphicFramePr",
				tag("a:graphicFrameLocks",
					attr("xmlns:a", "http://schemas.openxmlformats.org/drawingml/2006/main"),
					attr("noChangeAspect", "1"),
				),
			),
			tag("a:graphic",
				attr("xmlns:a", "http://schemas.openxmlformats.org/drawingml/2006/main"),
				tag("a:graphicData",
					attr("uri", "http://schemas.openxmlformats.org/drawingml/2006/picture"),
					tag("pic:pic",
						attr("xmlns:pic",
							"http://schemas.openxmlformats.org/drawingml/2006/picture",
						),
						// Non-Visual Image Properties
						// http://officeopenxml.com/drwPic-nvPicPr.php
						tag("pic:nvPicPr",
							tag("pic:cNvPr", attr("id", imgID), attr("name", img.name)),
							tag("pic:cNvPicPr"),
						),
						// Image Data
						// http://officeopenxml.com/drwPic-ImageData.php
						tag("pic:blipFill",
							tag("a:blip", attr("r:embed", relID), attr("cstate", "none")),
							tag("a:stretch", tag("a:fillRect")),
						),
						// Visual Image Properties
						// http://officeopenxml.com/drwSp-SpPr.php
						tag("pic:spPr",
							tag("a:xfrm",
								tag("a:off", attr("x", "0"), attr("y", "0")),
								tag("a:ext", attr("cx", cx), attr("cy", cy))),
							tag("a:prstGeom", attr("prst", "rect"),
								tag("a:avLst"),
							),
						),
					),
				),
			),
		),
	)
}

// nextImageID returns the next image id for this document
func (doc *Document) nextImageID() int {
	doc.lastImageIDOnce.Do(func() {
		nvPrs := doc.document.doc.FindElements("/" + "/pic:cNvPr")
		for _, nvPr := range nvPrs {
			id := nvPr.SelectAttr("id")
			if id == nil {
				continue
			}
			i, err := strconv.Atoi(id.Value)
			if err != nil {
				continue
			}
			if i > doc.lastImageID {
				doc.lastImageID = i
			}
		}
	})
	doc.lastImageID++
	return doc.lastImageID
}

// addImage adds an image at the current cursor
func (c *cursor) addImage(path string, x, y int) *cursor {
	switch c.e.token.(type) {
	case *body, *tableCell, *para, *listItem, *header, *link:
		c.startRun().addImage(path, x, y)
	case *run:
		img, err := newImg(c.doc, path, x, y)
		if err != nil {
			log.Error().Err(err).Msgf("cannot add img %s to document", path)
			return c
		}
		return c.AddXML(newImage(c.doc.nextImageID(), img))
	}
	return c
}
