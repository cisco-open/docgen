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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindParent(t *testing.T) {
	a := assert.New(t)
	r := newTestDocument().
		GetBookmark("main").
		AddChild(newParaEl(Props())).
		AddChild(newRunEl(Props()))

	a.NotNil(findParent[*run](r))
	a.NotNil(findParent[*para](r))
	a.NotNil(findParent[*body](r))
	a.Nil(findParent[*header](r))
	a.Equal(findParent[*run](r), r)
	a.Equal(findParent[*para](r), r.Parent())
	a.Equal(findParent[*body](r), r.Parent().Parent())
}
