# Using docgen as a Go Package

The `docgen` package provides functionality for generating Microsoft Word documents from templates with HTML content and template tags.

## Installation

```bash
go get github.com/ciscotools/docgen
```

## Quick Start

```go
package main

import (
    "log"
    "os"
    "strings"

    "github.com/ciscotools/docgen"
)

func main() {
    // Load a template document
    template, err := os.ReadFile("template.docx")
    if err != nil {
        log.Fatal(err)
    }

    // Create a new document from the template
    doc, err := docgen.NewDocument(template)
    if err != nil {
        log.Fatal(err)
    }

    // Insert HTML content at a bookmark
    htmlContent := strings.NewReader("<p>Hello, <b>World</b>!</p>")
    err = doc.InsertHTML(htmlContent, "main")
    if err != nil {
        log.Fatal(err)
    }

    // Render template tags
    ctx := map[string]interface{}{
        "title":  "My Document",
        "author": "John Doe",
        "date":   "2024-01-01",
    }
    err = doc.RenderTags(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Write the document to disk
    err = doc.WriteFile("output.docx")
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Document created successfully!")
}
```

## API Reference

### NewDocument

Creates a new Word document from a template.

```go
func NewDocument(template []byte) (*Document, error)
```

**Parameters:**
- `template`: Bytes of a valid .docx file (Word document in Office Open XML format)

**Returns:**
- `*Document`: A document that can be modified
- `error`: An error if the template is invalid or cannot be parsed

### Document Methods

#### WriteFile

Writes the document to a file.

```go
func (doc *Document) WriteFile(path string) error
```

**Parameters:**
- `path`: File path where the document will be written

**Returns:**
- `error`: An error if the file cannot be created or written

#### Write

Writes the document to an io.Writer.

```go
func (doc *Document) Write(w io.Writer) error
```

**Parameters:**
- `w`: Writer to receive the document bytes

**Returns:**
- `error`: An error if there are issues serializing or writing

#### InsertHTML

Inserts HTML content at a bookmark location.

```go
func (doc *Document) InsertHTML(input io.Reader, bookmark string) error
```

**Parameters:**
- `input`: Reader containing HTML content
- `bookmark`: Name of the bookmark where content will be inserted

**Returns:**
- `error`: An error if the bookmark is not found or HTML parsing fails

**Supported HTML elements:**
- Text formatting: `<b>`, `<i>`, `<u>`, `<s>`, `<sub>`, `<sup>`
- Paragraphs: `<p>`, `<div>`
- Lists: `<ul>`, `<ol>`, `<li>`
- Tables: `<table>`, `<tr>`, `<td>`, `<th>`
- Links: `<a href="...">`
- Images: `<img src="...">`
- Headings: `<h1>` through `<h6>`
- Line breaks: `<br>`

#### RenderTags

Renders template tags using the provided context.

```go
func (doc *Document) RenderTags(ctx any) error
```

**Parameters:**
- `ctx`: Template context (typically a map or struct) with values to substitute

**Returns:**
- `error`: An error if there are issues parsing or rendering tags

Template tags use the format `{{tagName}}` and can appear anywhere in the document (body, headers, footers).

#### GetBookmark

Finds and returns a bookmark element.

```go
func (doc *Document) GetBookmark(name string) *Element
```

**Parameters:**
- `name`: Name of the bookmark to find

**Returns:**
- `*Element`: Element representing the bookmark location, or nil if not found

## Creating Templates

Templates are standard Word documents (.docx) that include:

1. **Bookmarks**: Used as insertion points for HTML content
   - Create in Word: Insert > Bookmark
   - Use descriptive names like "main", "content", "header"

2. **Template Tags**: Placeholders in the format `{{tagName}}`
   - Can be placed anywhere in the document
   - Will be replaced with values from the context

## Examples

### Writing to HTTP Response

```go
func handler(w http.ResponseWriter, r *http.Request) {
    template, _ := os.ReadFile("template.docx")
    doc, _ := docgen.NewDocument(template)
    
    // Modify document...
    
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
    w.Header().Set("Content-Disposition", "attachment; filename=output.docx")
    doc.Write(w)
}
```

### Using with Buffers

```go
var buf bytes.Buffer
doc, _ := docgen.NewDocument(template)
doc.Write(&buf)

// buf now contains the complete .docx file
bytes := buf.Bytes()
```

## License

See [LICENSE](LICENSE) file for details.
