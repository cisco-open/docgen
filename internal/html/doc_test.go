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
package html

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterTemplatePath(t *testing.T) {
	a := assert.New(t)

	// Register a single path
	doc := New("main", nil)
	doc.RegisterTemplatePath("testdata")
	a.Contains(doc.Templates, "a1")
	a.Contains(doc.Templates, "a2")

	// Register multiple paths
	doc = New("main", nil)
	doc.RegisterTemplatePath("testdata/a")
	doc.RegisterTemplatePath("testdata/b")
	a.Contains(doc.Templates, "a1")
	a.Contains(doc.Templates, "a2")
	a.Contains(doc.Templates, "b1")
}

func TestRenderTemplate(t *testing.T) {
	a := assert.New(t)
	doc := New("main", nil)
	doc.RegisterTemplatePath("testdata")
	// Manually register a third template
	doc.RegisterTemplate("three", bytes.NewBuffer([]byte("three_body")))
	// Read two templates
	one := doc.Templates["a1"]
	two := doc.Templates["a2"]
	three := doc.Templates["three"]

	// Template one
	oneValue, err := one.Render(nil)
	a.NoError(err)
	a.Contains(oneValue, "one_body")

	// Template two
	twoValue, err := two.Render(nil)
	a.NoError(err)
	a.Contains(twoValue, "two_body")

	// Template three
	threeValue, err := three.Render(nil)
	a.NoError(err)
	a.Contains(threeValue, "three_body")
}

func TestIncludeHelper(t *testing.T) {
	a := assert.New(t)
	ctx := map[string]any{
		"root":  "root_value",
		"child": map[string]any{"value": "child_value"},
	}
	doc := New("main", ctx)
	// Validate that main context is used by default
	doc.RegisterTemplate("a", bytes.NewBuffer([]byte(`{{.root}} {{include "a_child"}}`)))
	doc.RegisterTemplate("a_child", bytes.NewBuffer([]byte(`{{.root}}`)))
	val, err := doc.Templates["a"].Render(ctx)
	a.NoError(err)
	a.Contains(val, "root_value root_value")
	a.NotContains(val, "child_value")

	// Validate that context is passed down when specified
	doc.RegisterTemplate("b", bytes.NewBuffer([]byte(`{{.root}} {{include "b_child" .child}}`)))
	doc.RegisterTemplate("b_child", bytes.NewBuffer([]byte(`{{.root}}{{.value}}`)))
	val, err = doc.Templates["b"].Render(ctx)
	a.NoError(err)
	a.Contains(val, "root_value child_value")

	// Validate child always has access to root
	doc.RegisterTemplate("c", bytes.NewBuffer([]byte(`{{.root}} {{include "c_child" .child}}`)))
	doc.RegisterTemplate("c_child", bytes.NewBuffer([]byte(`{{R.root}} {{.value}}`)))
	val, err = doc.Templates["c"].Render(ctx)
	a.NoError(err)
	a.Contains(val, "root_value root_value child_value")
}
