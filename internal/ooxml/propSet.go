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
	"strings"

	"github.com/beevik/etree"
)

type prop interface {
	addProp(*etree.Element)
}

// PropSet defines a unique set of properties
type PropSet struct {
	props []prop
	index map[prop]int
}

// Props is a shortcut to generate a PropSet from a list of properties
func Props(props ...prop) *PropSet {
	index := make(map[prop]int, len(props))
	for i, prop := range props {
		index[prop] = i
	}
	return &PropSet{props, index}
}

// Add adds a new property to a PropSet
func (set *PropSet) Add(v prop) *PropSet {
	if !set.Contains(v) {
		set.props = append(set.props, v)
		set.index[v] = len(set.props) - 1
	}
	return set
}

// Contains determines if a PropSet contains a given property
func (set *PropSet) Contains(v prop) bool {
	_, ok := set.index[v]
	return ok
}

// Difference determines set a - b
func (set *PropSet) Difference(other *PropSet) *PropSet {
	res := Props()
	for _, i := range set.Iter() {
		if !other.Contains(i) {
			res.Add(i)
		}
	}
	return res
}

// Empty return true if PropSet is empty
func (set *PropSet) Empty() bool {
	return len(set.props) == 0
}

// Equals determines if set a = b
func (set *PropSet) Equals(other *PropSet) bool {
	if len(set.props) != len(other.props) {
		return false
	}
	for _, i := range set.Iter() {
		if !other.Contains(i) {
			return false
		}
	}
	return true
}

// Iter creates an iterable list from a PropSet
func (set *PropSet) Iter() []prop {
	return set.props
}

func (set *PropSet) String() string {
	types := []string{}
	for _, i := range set.Iter() {
		types = append(types, fmt.Sprintf("%T", i))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(types, ", "))
}

// Union creates a new PropSet combing two existing PropSets
func (set *PropSet) Union(other *PropSet) *PropSet {
	res := Props()
	for _, i := range set.Iter() {
		res.Add(i)
	}
	for _, i := range other.Iter() {
		res.Add(i)
	}
	return res
}
