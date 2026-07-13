//go:build integration

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
package docgen_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cisco-open/docgen"
	"github.com/cisco-open/docgen/internal/templates"
	"github.com/stretchr/testify/assert"
)

// validateOOXML calls the Node.js validator to check if a .docx file is valid OOXML
func validateOOXML(t *testing.T, docxPath string) {
	t.Helper()

	// Get the path to the validator script
	validatorDir := filepath.Join("tools", "ooxml-validator-node")
	validatorScript := filepath.Join(validatorDir, "validate.js")

	// Check if validator exists
	if _, err := os.Stat(validatorScript); os.IsNotExist(err) {
		t.Fatalf("OOXML validator not found at %s. Run 'npm install' in %s", validatorScript, validatorDir)
	}

	// Run the validator
	cmd := exec.Command("node", validatorScript, docxPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("Validator output:\n%s", string(output))
		t.Fatalf("OOXML validation failed for %s: %v", docxPath, err)
	}

	t.Logf("OOXML validation passed for %s", filepath.Base(docxPath))
}

// TestIntegrationBasicDocument tests creating a basic document and validating it
func TestIntegrationBasicDocument(t *testing.T) {
	// Create a temporary file for the output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "basic.docx")

	// Create a document from the main template
	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err, "Failed to create document")

	// Insert basic HTML
	html := strings.NewReader("<p>This is a basic test document.</p>")
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err, "Failed to insert HTML")

	// Write the document
	err = doc.WriteFile(outputPath)
	assert.NoError(t, err, "Failed to write document")

	// Validate OOXML
	validateOOXML(t, outputPath)
}

// TestIntegrationFormattedText tests a document with various text formatting
func TestIntegrationFormattedText(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "formatted.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	html := strings.NewReader(`
		<h1>Main Heading</h1>
		<h2>Subheading</h2>
		<p>This paragraph contains <b>bold</b>, <i>italic</i>, <u>underline</u>, and <s>strikethrough</s> text.</p>
		<p>It also has <sub>subscript</sub> and <sup>superscript</sup> formatting.</p>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationLists tests documents with ordered and unordered lists
func TestIntegrationLists(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "lists.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	html := strings.NewReader(`
		<h2>Unordered List</h2>
		<ul>
			<li>First item</li>
			<li>Second item</li>
			<li>Third item</li>
		</ul>
		<h2>Ordered List</h2>
		<ol>
			<li>Step one</li>
			<li>Step two</li>
			<li>Step three</li>
		</ol>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationTable tests a document with a table
func TestIntegrationTable(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "table.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	html := strings.NewReader(`
		<h2>Device Table</h2>
		<table>
			<tr>
				<th>Device Name</th>
				<th>IP Address</th>
				<th>Status</th>
			</tr>
			<tr>
				<td>Router-01</td>
				<td>192.168.1.1</td>
				<td>Active</td>
			</tr>
			<tr>
				<td>Switch-01</td>
				<td>192.168.1.2</td>
				<td>Active</td>
			</tr>
			<tr>
				<td>Firewall-01</td>
				<td>192.168.1.3</td>
				<td>Standby</td>
			</tr>
		</table>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationLinks tests a document with hyperlinks
func TestIntegrationLinks(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "links.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	html := strings.NewReader(`
		<h2>Useful Links</h2>
		<p>Visit <a href="https://www.cisco.com">Cisco's website</a> for more information.</p>
		<p>Check the <a href="https://github.com/cisco-open/docgen">docgen repository</a> on GitHub.</p>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationComplexDocument tests a document with mixed content
func TestIntegrationComplexDocument(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "complex.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	html := strings.NewReader(`
		<h1>Network Assessment Report</h1>
		<p>This report provides an overview of the network infrastructure assessment conducted on <b>January 28, 2026</b>.</p>
		
		<h2>Executive Summary</h2>
		<p>The assessment covered the following areas:</p>
		<ul>
			<li>Network topology and design</li>
			<li>Device inventory and configuration</li>
			<li>Security posture</li>
			<li>Performance metrics</li>
		</ul>
		
		<h2>Device Inventory</h2>
		<table>
			<tr>
				<th>Device Type</th>
				<th>Model</th>
				<th>Quantity</th>
			</tr>
			<tr>
				<td>Router</td>
				<td>Cisco ISR 4000</td>
				<td>5</td>
			</tr>
			<tr>
				<td>Switch</td>
				<td>Catalyst 9300</td>
				<td>12</td>
			</tr>
			<tr>
				<td>Firewall</td>
				<td>ASA 5500-X</td>
				<td>2</td>
			</tr>
		</table>
		
		<h2>Recommendations</h2>
		<ol>
			<li>Upgrade firmware on all devices to the latest stable version</li>
			<li>Implement network segmentation using VLANs</li>
			<li>Enable logging and monitoring on all critical devices</li>
		</ol>
		
		<h2>References</h2>
		<p>For more information, visit <a href="https://www.cisco.com/c/en/us/support/">Cisco Support</a>.</p>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationWithTemplateTags tests a document with template tag rendering
func TestIntegrationWithTemplateTags(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "tags.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	// Render template tags first
	ctx := map[string]interface{}{
		"customer": "Example Corp",
		"date":     "January 28, 2026",
		"author":   "Network Team",
	}
	err = doc.RenderTags(ctx)
	assert.NoError(t, err)

	// Insert HTML content
	html := strings.NewReader(`
		<h1>Customer Report</h1>
		<p>This document was generated for our valued customer.</p>
	`)
	err = doc.InsertHTML(html, "main")
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}

// TestIntegrationEmptyTemplate tests creating a document without inserting content
func TestIntegrationEmptyTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "empty.docx")

	doc, err := docgen.NewDocument(templates.MainTemplate)
	assert.NoError(t, err)

	err = doc.WriteFile(outputPath)
	assert.NoError(t, err)

	validateOOXML(t, outputPath)
}
