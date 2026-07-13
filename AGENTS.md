# Copilot Instructions

## What This Project Does

`docgen` generates Microsoft Word (`.docx`) documents from HTML or Markdown
input files using a template `.docx` and a JSON context. It is both a CLI tool
(`cmd/docgen`) and an importable Go package (`github.com/cisco-open/docgen`).
There is also a Python wrapper in `python/`.

## Build, Test, and Lint

```bash
# Unit tests
go test ./...

# Run a single test
go test -run TestName ./internal/ooxml/

# Integration tests (require Node.js deps installed first)
cd tools/ooxml-validator-node && npm install && cd ../..
go test -tags=integration -v ./...

# Run a single integration test
go test -tags=integration -v -run TestIntegrationTable ./...
```

Integration tests write `.docx` files to a temp directory and validate them with
the Node.js OOXML validator at `tools/ooxml-validator-node/`.

## Architecture

```
docgen.go               Public Go package API (thin wrapper over internal/ooxml)
cmd/docgen/             CLI entrypoint (main.go, args.go, caf.go)
internal/ooxml/         Core OOXML engine — XML manipulation of .docx ZIP archives
internal/html/          HTML parsing helpers used by internal/ooxml
internal/templates/     Embedded .docx template files (MainTemplate, test template)
pkg/config/             YAML config file support for the CLI
python/                 Python wrapper library
tools/ooxml-validator-node/  Node.js OOXML validator used only in integration tests
```

**Data flow (CLI):** JSON context → render `{{tags}}` in template → parse
HTML/Markdown → convert to OOXML elements → insert at named bookmark → write
`.docx`.

**Data flow (package):** Caller loads a `.docx` template → calls
`docgen.NewDocument()` → calls `InsertHTML`/`InsertMarkdown` and `RenderTags` →
calls `WriteFile` or `Write`.

The public `docgen` package simply wraps `internal/ooxml.Document`. All real
work happens in `internal/ooxml`.

## Key Conventions

- **OOXML manipulation** uses `beevik/etree` for XML tree operations. The
  `.docx` ZIP structure is parsed in `internal/ooxml/zip.go`; the document parts
  (`document.xml`, `_rels/`, `numbering.xml`, headers, footers) are each held as
  an `ooxmlDoc` (an `etree.Document` + path).

- **Logging** uses `github.com/cisco-open/docgen/internal/log` (zerolog-based).
  Use `log.Debug()`, `log.Info()`, `log.Fatal()` — not `fmt` or the standard
  `log` package.

- **Error handling** uses `github.com/cockroachdb/errors`. Wrap errors with
  `errors.WithStack(err)` rather than `fmt.Errorf`.

- **Integration tests are gated** behind the `integration` build tag. All
  integration test functions are prefixed `TestIntegration` and live in
  `integration_test.go`. They follow this pattern: create doc from
  `templates.MainTemplate` → insert HTML → write to `t.TempDir()` → call
  `validateOOXML(t, path)`.

- **Default paragraph style** is `"Normal-6ptspacing"` (see
  `internal/ooxml/html.go`).

- **Input file routing**: `.md`/`.markdown` → Markdown path; `.html`/`.htm` →
  HTML path; anything else defaults to HTML (for backwards compatibility).

- **Template tags** (`{{tagName}}`) are handled separately from
  bookmark-inserted HTML. Tags appear anywhere in the document (body, headers,
  footers) and are replaced via `RenderTags`. The `Tag` struct in
  `internal/ooxml/tag.go` tracks the etree element positions for each tag across
  potentially split XML runs.

- **The `python/` directory** is an independent Python package with its own
  `pyproject.toml` and tests under `python/tests/`.

Refer to AGENTS.md at the repository root for all project instructions,
standards, and workflow guidelines.
