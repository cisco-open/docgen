# OOXML Validator Node

This directory contains a Node.js wrapper for validating Office Open XML (OOXML) documents using the `@xarsh/ooxml-validator` package.

## Purpose

This tool is used by the Go integration tests to validate that generated `.docx` files conform to the OOXML standard and are error-free.

## Installation

```bash
npm install
```

## Usage

```bash
node validate.js path/to/document.docx
```

The script will:
- Validate the document structure
- Print any validation errors found
- Exit with code 0 if valid, non-zero if invalid

## Integration with Go Tests

The Go integration tests (with `integration` build tag) call this validator after generating documents to ensure they are standards-compliant.

## Future

This Node-based validator provides immediate validation capability. In the future, a .NET OpenXML SDK validator may be added when it becomes available in Cisco IT's GitHub Actions environment.
