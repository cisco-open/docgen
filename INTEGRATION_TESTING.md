# Integration Testing

This document describes how to run OOXML integration tests for the docgen project.

## Overview

The integration tests validate that generated `.docx` files conform to the Office Open XML (OOXML) standard and are error-free. These tests:

1. Create documents from `templates.MainTemplate`
2. Insert various HTML test cases at the `main` bookmark
3. Write the resulting `.docx` to a temporary file
4. Call a Node.js validator to ensure the document is standards-compliant

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm (comes with Node.js)

## Running Integration Tests Locally

### 1. Install Node.js Dependencies

First, install the OOXML validator dependencies:

```bash
cd tools/ooxml-validator-node
npm install  # For local development
cd ../..
```

> **Note:** In CI/CD environments, the workflow uses `npm ci` for faster, more reliable builds. For local development, `npm install` is appropriate.

### 2. Run the Integration Tests

Run the tests with the `integration` build tag:

```bash
go test -tags=integration -v ./...
```

This will:
- Create test documents with various HTML content
- Validate each generated `.docx` file using the Node.js OOXML validator
- Report any validation errors

### 3. Test Individual Cases

You can run specific test cases using the `-run` flag:

```bash
# Run only the basic document test
go test -tags=integration -v -run TestIntegrationBasicDocument

# Run only table-related tests
go test -tags=integration -v -run TestIntegrationTable
```

## Available Test Cases

The integration test suite includes:

- **TestIntegrationBasicDocument** - Simple document with basic paragraph
- **TestIntegrationFormattedText** - Text with bold, italic, underline, etc.
- **TestIntegrationLists** - Ordered and unordered lists
- **TestIntegrationTable** - Document with tables
- **TestIntegrationLinks** - Document with hyperlinks
- **TestIntegrationComplexDocument** - Comprehensive document with mixed content
- **TestIntegrationWithTemplateTags** - Document with template tag rendering
- **TestIntegrationEmptyTemplate** - Empty template validation

## Continuous Integration

Integration tests run automatically in GitHub Actions:

- On every push to the `main` branch
- On every pull request to `main`

See `.github/workflows/ooxml-integration.yml` for the CI configuration.

## Validator Tool

The Node.js validator wrapper is located in `tools/ooxml-validator-node/`. It uses the `@xarsh/ooxml-validator` package to validate Office Open XML documents.

### Manual Validation

You can manually validate any `.docx` file:

```bash
cd tools/ooxml-validator-node
node validate.js path/to/document.docx
```

The validator will:
- Check the document structure
- Report any OOXML compliance errors
- Exit with code 0 if valid, non-zero if invalid

## Future Enhancements

This Node-based validation provides immediate OOXML validation capability. In the future, a .NET OpenXML SDK validator may be added when it becomes available in Cisco IT's GitHub Actions environment. The test structure is designed to accommodate additional validators alongside the Node validator.

## Troubleshooting

### Node.js Dependencies Not Found

If you see errors about missing Node.js modules:

```bash
cd tools/ooxml-validator-node
npm install
```

### Integration Tests Skipped

If integration tests are skipped when running `go test`, make sure to include the `-tags=integration` flag:

```bash
# Wrong - tests are skipped
go test -v ./...

# Correct - tests run
go test -tags=integration -v ./...
```

### Validation Errors

If a test fails with validation errors:

1. Check the test output for specific OOXML errors
2. Verify the HTML input is valid
3. Ensure the generated `.docx` file is being created correctly
4. You can inspect the generated `.docx` files in the temporary directory (shown in test output)

## Adding New Test Cases

To add new integration tests:

1. Open `integration_test.go`
2. Add a new test function with the `TestIntegration` prefix
3. Follow the pattern of existing tests:
   - Create a document from `templates.MainTemplate`
   - Insert HTML content
   - Write to a temporary file
   - Call `validateOOXML(t, outputPath)`

Example:

```go
func TestIntegrationNewCase(t *testing.T) {
    tmpDir := t.TempDir()
    outputPath := filepath.Join(tmpDir, "newcase.docx")

    doc, err := docgen.NewDocument(templates.MainTemplate)
    assert.NoError(t, err)

    html := strings.NewReader(`<p>Your HTML content here</p>`)
    err = doc.InsertHTML(html, "main")
    assert.NoError(t, err)

    err = doc.WriteFile(outputPath)
    assert.NoError(t, err)

    validateOOXML(t, outputPath)
}
```
