# docgen Python library

Python library for generating Microsoft Word documents using [docgen](https://github.com/ciscotools/docgen).

## Installation

```bash
pip install docgen
```

The `docgen` binary must also be available on `PATH`.
See the [docgen releases](https://github.com/ciscotools/docgen/releases) for installation instructions.

## Quick Start

```python
from docgen import generate

# Generate a document from an HTML file and a JSON context
generate(
    input="input.html",
    output="output.docx",
    context={"title": "My Document", "author": "Jane Doe"},
)
```

## API Reference

### `generate`

```python
def generate(
    input: Union[str, Path],
    output: Union[str, Path],
    context: Optional[Union[str, Path, dict]] = None,
    template: Optional[Union[str, Path]] = None,
    bookmark: str = "main",
    verbose: bool = False,
    workdir: Optional[Union[str, Path]] = None,
) -> None: ...
```

Generate a Word document using the `docgen` binary.

| Parameter  | Type                           | Description                                                   |
| ---------- | ------------------------------ | ------------------------------------------------------------- |
| `input`    | `str` or `Path`                | Path to the input file (`.html`, `.md`, or `.markdown`).     |
| `output`   | `str` or `Path`                | Path to the output `.docx` file.                             |
| `context`  | `str`, `Path`, `dict` or None  | Template variables – a JSON file path or a Python `dict`.    |
| `template` | `str`, `Path` or None          | Optional path to a custom `.docx` template file.             |
| `bookmark` | `str`                          | Bookmark name to render into (default: `"main"`).            |
| `verbose`  | `bool`                         | Enable verbose logging (default: `False`).                   |
| `workdir`  | `str`, `Path` or None          | Working directory for resolving relative paths.              |

Raises `DocgenError` if `docgen` cannot be found or exits with an error.

### `DocgenError`

Exception raised when `docgen` encounters an error.

## Development

```bash
cd python
pip install -e ".[dev]"
pytest
```
