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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindTagsSimple(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(newPara(newRun(newText("{{test}}"))))
	if a.Len(doc.tags, 1) {
		tag := doc.tags[0]
		a.Equal("test", tag.name)
		a.NotNil(tag.e1)
		a.NotNil(tag.e2)
		a.Equal(0, tag.i1)
		a.Equal(7, tag.i2)
	}
}

func TestFindTagsComplex(t *testing.T) {
	// More complex tag
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{.Test-test_test}}"))),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, ".Test-test_test")
	}
}

func TestFindTagsInvalid(t *testing.T) {
	// Invalid tag
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{te{test}}"))),
	)
	a.Len(doc.tags, 0)

	doc.updateDocument(
		newPara(newRun(newText("{{test}"))),
	)
	a.Len(doc.tags, 0)

	doc.updateDocument(
		newPara(newRun(newText("{test}}"))),
	)
	a.Len(doc.tags, 0)
}

func TestFindTagsTwoRuns(t *testing.T) {
	// Find tag broken between runs
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("{{te")),
			newRun(newText("st}}")),
		),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, "test")
	}
}

func TestFindTagsThreeRuns(t *testing.T) {
	// Tag *really* broken between runs
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("{{t")),
			newRun(newText("es")),
			newRun(newText("t}}")),
		),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, "test")
	}
}

func TestFindTagsExtraBrackets(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{{test}}"))),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, "test")
	}
}

func TestFindTagsExtraneousBrackets(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{before{{test}}after}}"))),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, "test")
	}
}

func TestFindTagsOpeningBracketSplit(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{"), newText("{test}}"))),
	)
	if a.Len(doc.tags, 1) {
		a.Equal(doc.tags[0].name, "test")
	}
}

func TestFindTagsMultiple(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{one}}"), newText("{{two}}"))),
	)
	if a.Len(doc.tags, 2) {
		a.Equal(doc.tags[0].name, "one")
		a.Equal(doc.tags[1].name, "two")
	}
}

func TestFindTagsParaReset(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		tag("w:body",
			newPara(
				newRun(newText("{{te")),
			),
			newPara(
				newRun(newText("st}}")),
			),
		),
	)
	a.Len(doc.tags, 0)
}

func TestFindTagsMultipleParas(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		tag("w:body",
			newPara(newRun(
				newText("{{one}}")),
				newText("{{two}}")),
			newPara(newRun(
				newText("{{one}}")),
			),
			newPara(newRun(
				newText("{{two}}")),
			),
		),
	)
	a.Len(doc.tags, 4)
}

func TestTagReplaceSimple(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{test}}"))),
	)
	doc.tags[0].render("content")
	expected := newOOXMLDocFromTag(
		newPara(newRun(newText("content"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	// Test case for GitHub issue: Tag rendering with multiple tags in same element
	// When rendering {{.Report}} with "ACI EoX Report" followed by another tag,
	// the indices of subsequent tags can become stale
	doc.updateDocument(
		newPara(newRun(newText("{{.Report}}{{.Other}}"))),
	)
	// Use RenderTags to test the fix for rendering multiple tags in same element
	err := doc.RenderTags(map[string]any{
		"Report": "ACI EoX Report",
		"Other":  " Document",
	})
	a.NoError(err)
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("ACI EoX Report Document"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagReplaceSurrounding(t *testing.T) {
	// Surrounding content
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("before|{{test}}|after"))),
	)
	doc.tags[0].render("content")
	expected := newOOXMLDocFromTag(
		newPara(newRun(newText("before|content|after"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	doc.updateDocument(
		newPara(newRun(newText("{{test}}|after"))),
	)
	doc.tags[0].render("content")
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("content|after"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	doc.updateDocument(
		newPara(newRun(newText("before|{{test}}"))),
	)
	doc.tags[0].render("content")
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("before|content"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	// Test case for issue: Tag rendering with replacement text that contains tag name suffix
	// This reproduces the bug where "{{.Report}}" replaced with "ACI EoX Report" 
	// incorrectly produces "ACI EoX Reportort"
	doc.updateDocument(
		newPara(newRun(newText("{{.Report}}"))),
	)
	doc.tags[0].render("ACI EoX Report")
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("ACI EoX Report"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	// Test with prefix text that contains part of the replacement
	doc.updateDocument(
		newPara(newRun(newText("Rep{{.Report}}"))),
	)
	doc.tags[0].render("ACI EoX Report")
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("RepACI EoX Report"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())

	// Test case: Two tags in the same element, rendering via RenderTags
	doc.updateDocument(
		newPara(newRun(newText("{{.tag1}}{{.tag2}}"))),
	)
	a.Len(doc.tags, 2, "Should find 2 tags")
	// Use RenderTags to render both tags (which will handle them in reverse order)
	err := doc.RenderTags(map[string]any{
		"tag1": "AAAA",
		"tag2": "BBBB",
	})
	a.NoError(err)
	expected = newOOXMLDocFromTag(
		newPara(newRun(newText("AAAABBBB"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagReplaceTwoTextTags(t *testing.T) {
	// Split across two text elements in a single run
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(
				newText("before|{{te"),
				newText("st}}|after"),
			),
		),
	)
	doc.tags[0].render("content")
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(
				newText("before|content"),
				newText("|after"),
			),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagReplaceThreeTextTags(t *testing.T) {
	// Split across three text elements in a single run
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(
				newText("before|{{t"),
				newText("es"),
				newText("t}}|after"),
			),
		),
	)
	doc.tags[0].render("content")

	// This avoids creating a self-closing tag which emulates the modified test tag
	emptyTxt := newText(" ")
	emptyTxt.SetText("")
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(
				newText("before|content"),
				emptyTxt,
				newText("|after"),
			),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagReplaceTwoRuns(t *testing.T) {
	// Split across two runs
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("{{te")),
			newRun(newText("st}}")),
		),
	)
	doc.tags[0].render("content")

	emptyTxt := newText(" ")
	emptyTxt.SetText("")
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("content")),
			newRun(emptyTxt),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagReplaceThreeRuns(t *testing.T) {
	// Split across three runs
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("before|{{t")),
			newRun(newText("es")),
			newRun(newText("t}}|after")),
		),
	)
	doc.tags[0].render("content")

	emptyTxt := newText(" ")
	emptyTxt.SetText("")
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("before|content")),
			newRun(emptyTxt),
			newRun(newText("|after"))))
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestRenderTags(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("before|{{.tag}}|after")),
		),
	)
	doc.RenderTags(map[string]any{
		"tag": "content",
	})
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("before|content|after"))))
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestRenderTagsWithHelper(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("before|{{len .tag}}|after")),
		),
	)
	doc.RenderTags(map[string]any{
		"tag": []int{1, 2, 3},
	})
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("before|3|after"))))
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestRenderTagsWithSprigHelper(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("before|{{first .tag}}|after")),
		),
	)
	doc.RenderTags(map[string]any{
		"tag": []int{1, 2, 3},
	})
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText("before|1|after")),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestRenderTagsWithDate(t *testing.T) {
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(
			newRun(newText("{{date}}")),
		),
	)
	doc.RenderTags(map[string]any{
		"tag": []int{1, 2, 3},
	})
	expected := newOOXMLDocFromTag(
		newPara(
			newRun(newText(time.Now().Format("2 Jan 2006"))),
		),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagRenderWithStaleIndices(t *testing.T) {
	// Regression test for GitHub issue: Tag rendering duplicating content
	// This test simulates the bug where render() is called with stale indices
	// after the text has been modified by another operation
	a := assert.New(t)
	doc := newTestDocument()

	// Create a paragraph with the tag "{{.Report}}"
	doc.updateDocument(
		newPara(newRun(newText("{{.Report}}"))),
	)

	// We should have found one tag
	a.Len(doc.tags, 1)
	tag := doc.tags[0]

	// Manually modify the text to "ACI Healthcheck" to simulate what happens
	// when another tag is rendered first or some other operation modifies the text
	tag.e1.SetText("ACI Healthcheck")

	// Now try to render the tag with stale indices
	// The bug was that render() would read "heck" as suffix from the modified text
	// and produce "ACI Healthcheckheck" instead of "ACI Healthcheck"
	tag.render("ACI Healthcheck")

	// Check the result - should be "ACI Healthcheck", not "ACI Healthcheckheck"
	expected := newOOXMLDocFromTag(
		newPara(newRun(newText("ACI Healthcheck"))),
	)
	result := doc.document.StripBody().ToString()
	a.Equal(expected.ToString(), result, "Tag rendering with stale indices should not duplicate content")
}

func TestTagRenderingDuplication(t *testing.T) {
	// Regression test for GitHub issue: Tag rendering duplicating content
	// When rendering {{.Report}} with "ACI Healthcheck", it was producing "ACI Healthcheckheck"
	a := assert.New(t)
	doc := newTestDocument()
	doc.updateDocument(
		newPara(newRun(newText("{{.Report}}"))),
	)
	err := doc.RenderTags(map[string]any{
		"Report": "ACI Healthcheck",
	})
	a.NoError(err)
	expected := newOOXMLDocFromTag(
		newPara(newRun(newText("ACI Healthcheck"))),
	)
	a.Equal(expected.ToString(), doc.document.StripBody().ToString())
}

func TestTagRenderWithStaleIndicesMultiElement(t *testing.T) {
	// Test for tag spanning multiple text elements within same run
	a := assert.New(t)
	doc := newTestDocument()

	// Create a tag split across two text elements
	doc.updateDocument(
		newPara(newRun(
			newText("{{.Rep"),
			newText("ort}}"))),
	)

	a.Len(doc.tags, 1)
	tag := doc.tags[0]

	// Modify the first element's text
	tag.e1.SetText("ACI Health")

	// Modify the second element's text
	tag.e2.SetText("check")

	// Now render - should handle stale indices gracefully
	tag.render("ACI Healthcheck")

	// Verify the exact output structure
	// For multi-element tags, the implementation sets e1 to "pfx + txt" and e2 to "sfx"
	// With modified text, pfx and sfx should be empty, so:
	// e1 should be "ACI Healthcheck" and e2 should be ""
	result := doc.document.StripBody().ToString()
	a.NotContains(result, "ACI Healthcheckheck", "Should not duplicate content")
	a.Contains(result, "ACI Healthcheck", "Should contain the replacement text")

	// Verify that the text elements are correctly set
	a.Equal("ACI Healthcheck", tag.e1.Text(), "e1 should contain the replacement text")
	a.Equal("", tag.e2.Text(), "e2 should be empty (no suffix)")
}
