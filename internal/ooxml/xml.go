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

// XML util

import (
	"github.com/cisco-open/docgen/internal/log"

	"github.com/beevik/etree"
)

func attr(key, value string) *etree.Attr {
	return &etree.Attr{Key: key, Value: value}
}

func tag(tag string, mods ...any) *etree.Element {
	e := etree.NewElement(tag)
	for _, mod := range mods {
		switch m := mod.(type) {
		case *etree.Element:
			e.AddChild(m)
		case *etree.Attr:
			e.CreateAttr(m.Key, m.Value)
		case *etree.CharData:
			e.CreateText(m.Data)
		case string:
			e.CreateText(m)
		case []etree.Token:
			for _, child := range m {
				e.AddChild(child)
			}
		case []*etree.Element:
			for _, child := range m {
				e.AddChild(child)
			}
		default:
			log.Warn().Msgf("unexpected tag mod type:%T for %s", mod, tag)
		}
	}
	return e
}
