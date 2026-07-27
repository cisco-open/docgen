## Overview

Docgen automates generating _high-quality_ documents using declarative HTML or
Markdown templates.

See the [Quick Start](#quick-start) to get started immediately.

---

This tool consumes the output of your existing automation tools, e.g. analysis
tools, build tools, etc., and generates finished documentation from this
output. It does this accurately, efficiently, and reliably, and generates
high-quality documentation.

See the [Comparison](#comparison-to-other-solutions) section for comparison with
other solutions, which will further illustrate the unique value this tool
provides.

## Quick Start

### Configuration Options

Docgen supports two ways to configure its behavior:

1. **Command-line arguments** (default): Pass parameters directly via CLI
2. **YAML configuration file**: Use a config file for easier management and
   reusability

#### Using Command-Line Arguments

The traditional way to use docgen is with command-line arguments:

```shell
docgen --input input.html --output out.docx --context context.json --verbose
```

Available command-line options:

```shell
docgen --help
```

#### Using a Configuration File

For easier management and reusability, you can use a YAML configuration file:

```shell
docgen --config config.yaml
```

**Important:** When a config file is provided, all command-line arguments
(except `--config`) are ignored.

##### Example Configuration File

Create a file named `config.yaml`:

```yaml
# Path to the tag context JSON file
context: context.json

# Path to the input file (.html, .md, or .markdown)
input: input.html

# Path to the output Word document file
output: out.docx

# Path to a custom Word document template (optional)
template: ""

# Bookmark name to render content into
bookmark: main

# Path to the log file
logfile: docgen.log

# Enable verbose (debug level) logging
verbose: false

# Working directory for relative paths
workdir: ""
```

See [config-example.yaml](config-example.yaml) for a complete example with
documentation.

##### Benefits of Using Configuration Files

- **Reusability**: Save different configurations for different projects
- **Version control**: Check configuration files into git alongside your
  templates
- **Simplicity**: Avoid long command lines with many arguments
- **Documentation**: YAML comments can document your configuration choices

### Simple Example

Download the latest docgen release from
[GitHub releases](https://github.com/cisco-open/docgen/releases).

Clone the example repo:

```shell
git clone --depth 1 https://github.com/cisco-open/docgen-example.git
```

Copy or move the `docgen` (`docgen.exe` for Windows) file into your
`docgen-example` folder that you just cloned.

Run the development tool while in that folder: `./docgen --dev`. This starts the
development environment, listens for changes to the files, validates your
templates for correctness, and opens a document preview in your default browser.

#### Understanding the Simple Example

Two things are provided in this example:

1. An example `input.html` file. This is the HTML content that will be inserted
   into the deliverable document. In a real workflow, you generate this file
   from your existing automation tool using a templating engine suited to your
   language — e.g. Jinja2 for Python, Handlebars for JavaScript, or Templ for
   Go. See the [Generating Input Content](#generating-input-content) section for
   examples.
2. Optionally, a `context.json` file for populating `{{tag}}` placeholders
   directly in the Word document template (e.g. cover page fields, headers,
   footers). See [Direct Template Tags](#direct-template-tags).

Running the development server does the following:

1. Reads the `input.html` (or `.md`) file as the document body content.
2. Inserts that content into the `main` bookmark in the output document.
   Docgen uses a built-in default document template, but you can supply
   your own custom template with pre-existing content.
3. Optionally renders `{{tag}}` placeholders in the Word template using values
   from `context.json`.

See the [Comprehensive Guide](#comprehensive-guide) for detailed usage and
examples.

## Comprehensive Guide

### Generating Input Content

Docgen accepts a pre-rendered HTML or Markdown file as its input. You generate
this file using whatever templating solution is idiomatic for your tool's
language. This keeps docgen focused on document conversion while letting you use
mature, well-documented templating engines you already know.

#### Python — Jinja2

[Jinja2](https://jinja.palletsprojects.com/) is the standard templating engine
for Python tools.

```python
from jinja2 import Environment, FileSystemLoader

env = Environment(loader=FileSystemLoader("."))
template = env.get_template("report.html.j2")

data = {
    "devices": [
        {"name": "switch01", "ip": "10.0.0.1"},
        {"name": "switch02", "ip": "10.0.0.2"},
    ]
}

with open("input.html", "w") as f:
    f.write(template.render(**data))
```

**`report.html.j2`**

```html
<h1>Device Report</h1>
<table>
  <tr>
    <th>Name</th>
    <th>IP</th>
  </tr>
  {% for device in devices %}
  <tr>
    <td>{{ device.name }}</td>
    <td>{{ device.ip }}</td>
  </tr>
  {% endfor %}
</table>
```

Alternatively, [FastHTML](https://fastht.ml/) can generate HTML directly from
Python objects without a separate template file.

#### JavaScript / Node.js — Handlebars

[Handlebars](https://handlebarsjs.com/) is a popular choice for JavaScript
tools.

```javascript
const Handlebars = require("handlebars");
const fs = require("fs");

const template = Handlebars.compile(fs.readFileSync("report.html.hbs", "utf8"));

const data = {
  devices: [
    { name: "switch01", ip: "10.0.0.1" },
    { name: "switch02", ip: "10.0.0.2" },
  ],
};

fs.writeFileSync("input.html", template(data));
```

**`report.html.hbs`**

```html
<h1>Device Report</h1>
<table>
  <tr>
    <th>Name</th>
    <th>IP</th>
  </tr>
  {{#each devices}}
  <tr>
    <td>{{name}}</td>
    <td>{{ip}}</td>
  </tr>
  {{/each}}
</table>
```

#### Go — Templ

[Templ](https://templ.guide/) is a type-safe HTML templating library for Go, and
is the recommended approach for Go-based tools.

```go
// devices.templ
package report

templ DeviceTable(devices []Device) {
    <h1>Device Report</h1>
    <table>
        <tr><th>Name</th><th>IP</th></tr>
        for _, d := range devices {
            <tr><td>{ d.Name }</td><td>{ d.IP }</td></tr>
        }
    </table>
}
```

```go
// main.go
f, _ := os.Create("input.html")
defer f.Close()
report.DeviceTable(devices).Render(context.Background(), f)
```

Go's standard `html/template` or `text/template` packages also work if you
prefer not to add a dependency.

#### Markdown

If your content is simple prose with basic tables or lists, generating Markdown
directly (e.g. with string formatting or a Markdown builder library) is often
the simplest approach. Docgen accepts `.md` and `.markdown` files in addition to
`.html`.

### Direct Template Tags

In addition to inserting content from an input file, you can embed `{{tag}}`
placeholders directly in your Word document template. These are rendered using a
JSON context file and are useful for fields that appear in the document template
itself — such as cover page fields, headers, or footers — rather than in the
body content.

Pass a JSON context file with `--context context.json`:

```json
{
  "customer": "Acme Corp",
  "date": "26 Jun 2026",
  "author": "Jane Smith"
}
```

Any `{{customer}}`, `{{date}}`, or `{{author}}` tags in the Word template are
replaced with the corresponding values. These tags use Go's `text/template`
syntax, so simple expressions and conditionals are supported, but complex logic
belongs in your input generation step.

#### Advanced HTML Examples

##### Column widths

Table columns widths can be specified with the `colgroup` HTML element:

```html
<table>
  <colgroup>
    <col style="width: 30%" />
    <col style="width: 70%" />
  </colgroup>
  <!--rest of table-->
</table>
```

All tables will be rendered at full page width. Without `colgroup` columns will
be evenly distributed, e.g. three columns will render at 33% each.

##### Rowspan and Colspan

`rowspan` and `colspan` allow merging table cells horizontally and vertically
respectively. The
[Mozilla development docs](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td)
explain these fields.

`colspan` is pretty straightforward. The following example illustrates a table
with _three_ columns, where two of the header cells are merged into a single
header:

```html
<table>
  <tr>
    <th>Name</th>
    <th colspan="2">Address</th>
  </tr>
  <tr>
    <td>switch01</td>
    <td>10.0.0.1</td>
    <td>/24</td>
  </tr>
</table>
```

`rowspan` offers a unique challenge. Namely, the spanned cell is only rendered
in the first row (the one with the rowspan attribute). From there on, that cell
is excluded in the HTML. This is an example of correctly formatted HTML using
rowspan:

```html
<table>
  <tr>
    <th>Model</th>
    <th>Switches</th>
  </tr>
  <tr>
    <th rowspan="2">93180YC-EX</th>
    <th>switch01</th>
  </tr>
  <tr>
    <th>switch02</th>
  </tr>
</table>
```

Note, specifically, that the spanned row indicating the model is excluded from
the second row. That's because this field is spanned across two rows. When
generating this HTML from your templating engine, you need to emit the `rowspan`
cell only in the first row of each group and omit it from subsequent rows.

Here's how that looks in Jinja2:

```html
<table>
  <tr>
    <th>Model</th>
    <th>Switch</th>
  </tr>
  {% for model, switches in devices_by_model.items() %} {% for switch in
  switches %}
  <tr>
    {% if loop.first %}
    <td rowspan="{{ switches|length }}">{{ model }}</td>
    {% endif %}
    <td>{{ switch.name }}</td>
  </tr>
  {% endfor %} {% endfor %}
</table>
```

##### Images

Images are supported with the HTML image tag:

```html
<img src="/router.png" />
```

Images must be in either `png` or `jpeg` format. The preceeding slash will be
read by the development server, which will try to read the image in relation to
the current working directory. In this example, `router.png` is expected to
exist in the current folder.

## Comparison to other solutions

This solution was created because of a distinct need for better document
automation. The following solutions were considered, and have all been used in
production leading up to this tool:

**Programmtic document creation** is where you (the tool owner) use a docx
library to progressively build the document throughout your code. For example,
the [python-docx library](https://python-docx.readthedocs.io/en/latest/)
provides method calls for adding paragraphs, tables, headers, formmating, etc.
This is a very common solution and is how most existing automated documentation
creation is performed.

**Challenges:**

- Mixing code logic with presentation logic makes it very difficult to write and
  maintain quality documentation. Programming languges such as Python were not
  designed for writing visual documentation.

**Pandoc** is a universal document transformation library which will translate
from Markdown or HTML to docx. Other solutions can be used to create an HTML or
Markdown document (e.g. Jinja2, Handlebars, Templ, etc), and then Pandoc can be
used to transform this into a Microsoft Word document. This was the predecesor
to docgen and is the closest in architecture.

**Challenges:**

- Getting started, developing, and deploying a tool that uses this solution
  requires a lot more dependencies.
- Invalid Markdown or HTML can produce a broken/invalid Word document.
- Pandoc is not as fast or efficient as a purpose-built tool such as docgen.

### Feature Comparison

| Feature                                       | docgen             | Library            | Pandoc             |
| --------------------------------------------- | :----------------- | :----------------- | :----------------- |
| Automated                                     | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Separation of business and presentation logic | :white_check_mark: | :x:                | :white_check_mark: |
| Declarative                                   | :white_check_mark: | :x:                | :white_check_mark: |
| Easy to integrate with new tools              | :white_check_mark: | :x: **¹**          | :x: **²**          |
| Git-friendly                                  | :white_check_mark: | :x: **³**          | :white_check_mark: |
| Will not create a broken Word document        | :white_check_mark: | :white_check_mark: | :x:                |
| Fast / memory efficient                       | :white_check_mark: | :white_check_mark: | :x:                |
| Developer Friendly                            | :white_check_mark: | :white_check_mark: | :x: **⁴**          |

**¹** A docx library such as python-docx can be very easy to get started with,
but writing _quality_ documentation in Python can be very challenging.

**²** Several independent components are required for this solution, making it
much more difficult for an end-user to run.

**³** Mixing business and presentation logic makes it very difficult to
determine which commits changes are documentation changes vs core code changes.
It also makes it more difficult to wind these changes back independently.

**⁴** Templating + pandoc + a web server + live reload requires installing
multiple components or running the development environment in docker, which
significantly increases the "getting started" burden for a new developer.

## Contributing

This is a community tool. Contributions, bug fixes, and feature requests are
welcome.

### Testing

#### Unit Tests

Run the standard unit tests:

```bash
go test ./...
```

#### Integration Tests

Integration tests validate that generated `.docx` files conform to the OOXML
standard. See [INTEGRATION_TESTING.md](INTEGRATION_TESTING.md) for detailed
instructions on running integration tests locally.

Quick start:

```bash
# Install Node.js dependencies (local development)
cd tools/ooxml-validator-node
npm install
cd ../..

# Run integration tests
go test -tags=integration -v ./...
```

### Architecture

#### Architectural Overview

#### Deep Dive

### Coding Standards

## License

SPDX-License-Identifier: Apache-2.0

Copyright 2026 Cisco Systems, Inc. and their affiliates

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
