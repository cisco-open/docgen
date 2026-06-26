## Overview

Docgen automates generating _high-quality_ CX deliverable documents from our
tools.

See the [Quick Start](#quick-start) to get started immediately.

---

Writing our deliverable documents is often a time-consuming and error-prone
process. Even with good templates, e.g. in SCDP, populating customer-specific
information can be a consderable effort. In CX, we've done a great job with
automating other aspects of delivery, but deliverable document generation is
still often harder than it should be.

This tool consumes the output of our existing automation tools, e.g. analysis
tools, build tools, etc., and generates deliverable documentation from this
output. It does this accurately, efficiently, and reliably, and generates
high-quality documentation that conforms to our CX templating standards.

As a tool owner, you can integrate this tool into your existing tool to bolt on
document generation or improve on existing, difficult to maintain solutions. As
an end-user doing CX delivery, this tool will transparently generate your
deliverable documents providing significant time and cost savings delivery.

See the [Comparison](#comparison-to-other-solutions) section for comparison with
other solutions, which will further illustrate the unique value this tool
provides.

## Quick Start

### Configuration Options

Docgen supports two ways to configure its behavior:

1. **Command-line arguments** (default): Pass parameters directly via CLI
2. **YAML configuration file**: Use a config file for easier management and reusability

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

**Important:** When a config file is provided, all command-line arguments (except `--config`) are ignored.

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

See [config-example.yaml](config-example.yaml) for a complete example with documentation.

##### Benefits of Using Configuration Files

- **Reusability**: Save different configurations for different projects
- **Version control**: Check configuration files into git alongside your templates
- **Simplicity**: Avoid long command lines with many arguments
- **Documentation**: YAML comments can document your configuration choices

### Simple Example

Download the latest docgen release from
[GitHub releases](https://github.com/ciscotools/docgen/releases).

Clone the example repo:

```shell
git clone --depth 1 https://github.com/ciscotools/docgen-example.git
```

Copy or move the `docgen` (`docgen.exe` for Windows) file into your
`docgen-example` folder that you just cloned.

Run the development tool while in that folder: `./docgen --dev`. This starts the
development environment, listens for changes to the files, validates your
templates for correctness, and opens a document preview in your default browser.

#### Understanding the Simple Example

Two things are provided in this example:

1. An example JSON `input.json` file. This will the output from your existing
   automation tool, e.g. API output, health check, day-1 build tool, etc. If
   your tool doesn't currently produce JSON, it probably can. See the
   [JSON Generation Examples](#json-generation-examples) section for ideas.
2. Example templates. Templates are written in HTML and use a simple templating
   language for basic logic. Detailes can be found in the
   [Templating](#templating) section of this guide and examples in the
   [Template Examples](#template-examples) section. The simple example provides
   examples of common patterns that will be present in _most_ templates.

Running the development server does the following:

1. Reads `input.json` as an input "context"
2. Uses the context to render the HTML templates. The starting template is
   `main.html` and other templates are included with the `include` statement.
   This generates a single `out.html` file representing your deliverable
   document.
3. The `out.html` content gets rendered into the `main` bookmark in the
   deliverable document. Docgen is using a built-in CX document template for
   this, but you can also use your own custom template containing pre-existing
   content.

See the [Comprehensive Guide](#comprehensive-guide) for detailed usage and
examples.

## Comprehensive Guide

### JSON input

The input data is standard JSON.

#### Basic Python

If your tool is written in Python, data can be written to JSON directly, using
the built-in `json` library:

```python
import json
my_app_data = {"devices": [{"name": "device01"}, {"name": "device02"}]}
with open("input.json") as f:
  f.write(json.dumps(my_app_data))
```

If you're programmatically generating a document, e.g. with `python-docx` or
similar, you can instead write your data to a single dictionary and dump it to
JSON.

#### CSV

In Python, CSV can be read from a file with the built-in CSV module. This
example assumes your CSV has a header row and uses the headers as the dictionary
keys:

```python
import csv

with open("data.csv") as f:
  csv_reader = csv.DictReader(f, delimiter=',')
with open("input.json") as f:
  f.write(json.dumps(list(csv_reader)))
```

This module also works with built in objects as opposed to reading from a file:

```python
import csv

data = ["name,ip", "device01,10.0.0.1", "device02,10.0.0.2"]
csv_reader = csv.DictReader(["header1"], delimiter=',')
with open("input.json") as f:
  f.write(json.dumps(list(csv_reader)))

# [{"name": "device01", "ip": "10.0.0.1"}, {"name": "device02", "ip": "10.0.0.2"}]
```

#### Excel

If your data is in Excel, you can either save it to CSV or use the `openpyxl`
library to parse it into Python objects. Once in Python, use the built in JSON
module to write to JSON.

### Templating

Templating is performed by Go's text/template library. Basic documentation on
this library is provided in the
[Go documentation for text/template](https://pkg.go.dev/text/template).

Additionally, the [sprig helper functions](http://masterminds.github.io/sprig/)
are included, which provide a number of useful helper functions.

If coming from Jinja2, you may notice these templates have very limited features
in comparison. Specifically, complex code is not written into templates. This is
intentional and is a feature, as it encourages the separation of code logic from
view logic. This separation is an important principle of maintainable code.

#### Basic building blocks

##### Variables

Variables are specified with a dot notation, e.g.:

```json
{ "name": "device01" }
```

```html
<h1>Configuration for {{.name}}</h1>
```

##### `if` statement

The `if` block conditionally renders its contents. For example:

```json
{ "is_healthy": true }
```

```html
{{if .is_healthy}}
<p>The network is healthy.</p>
{{end}}
```

result:

```html
<p>The network is healthy.</p>
```

If statements can also include an else clause:

```json
{ "is_healthy": false }
```

```html
{{if .is_healthy}}
<p>The network is healthy.</p>
{{else}}
<p>Unhealthy!</p>
{{end}}
```

result:

```html
<p>Unhealthy!</p>
```

`if` statements evaluate "truthy" values and not just booleans. This means `0`,
`""`, `[]`, or even a missing key will evaluate to false.

Note that the `docgen --dev` mode warns you of missing keys on the CLI. This
validation checks that keys mentioned in the template exist within the JSON and
is to help you avoid mistyping key names. However, missing keys are still valid
templating code and will still render correctly. If you want to turn off the
template validation, run `docgen --dev --disable-validation`.

##### `with` statement

The `with` statement is similar to if, but changes the context to match the
variable specified in the with statement.

Suppose we want to conditionally render some text based on chassis health.
Here's our example `input.json`:

```json
{
  "chassis": {
    "name": "device01",
    "is_healthy": false
  },
  "module_count": 3
}
```

This is how it might look while **not** using the `with` statement:

```html
{{if .chassis.is_healthy}}
<p>{{.chassis.name}} is healthy!</p>
{{else}}
<p>{{.chassis.name}} is unhealthy!</p>
{{end}}
```

This is how it would look using the `with` statement. Notice that we shifted
_into_ the `.chassis` context, so no longer need to refer to `.chassis` every
time. This is an overly simplified example, but this can be very powerful in
more complex templates, e.g. where the outer complex is referenced repeatedly.

```html
{{with .chassis}} {{if .is_healthy}}
<p>{{.name}} is healthy!</p>
{{else}}
<p>{{.name}} is unhealthy!</p>
{{end}} {{end}}
```

When using the `with` statement, you can always refer to the templating root
with the special `$` indicator. Suppose in the previous example, we also want to
refer to the `.modules_count` field. Note that this field is not under
`.chassis`.

```html
{{with .chassis}} {{if .is_healthy}}
<p>{{.name}} is healthy!</p>
{{else}}
<p>{{.name}} is unhealthy!</p>
{{end}}
<p>{{.module_count}} modules found.</p>
{{end}}
```

##### `range` statement

The `range` block loops over a list or object (dictionary):

```json
{
  "devices": [
    { "name": "switch01" },
    { "name": "switch02" },
    { "name": "switch03" }
  ]
}
```

```html
<ul>
  {{range .devices}}
  <li>{{.name}}</li>
  {{end}}
</ul>
```

You can also assign the items to variables. Note that this time, devices is an
obect instead of an array.

```html
<ul>
  {{range $index, $device := .devices}}
  <li>{{$device.name}} is device number {{$index}}.</li>
  {{end}}
</ul>
```

##### `include` helper

The include helper is a custom helper that lets you include additional
templates. For example:

**`main.html`**

```html
<h1>Devices</h1>
<p>This section outlines the device details for {{customer}}.</p>
{{include "devices"}}
```

**`devices.html`**

```html
{{if .devices}}
<p>The following devices were found in your network:</p>
<ul>
  {{range .devices}}
  <li>{{.}}</li>
  {{end}}
</ul>
{{else}}
<p>No devices found.</p>
{{end}}
```

##### `date` helper

The `date` helper inserts the current date in RFC53222 format, e.g. 21 Jan 2013.
This format is generaly more user friendly than ISO year-first formats while
avoiding US/Europe ambiguity of MM/DD vs DD/MM.

##### External documentation

- [Hashicorp's Nomad documentation](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax)
- [Go text/template documentation](https://pkg.go.dev/text/template)
- [Sprig helper functions](http://masterminds.github.io/sprig/)

#### Advanced Template Examples

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
  {{range .devices}}
  <tr>
    <td>{{.name}}</td>
    <td>{{.ip}}</td>
    <td>{{.pfx}}</td>
  </tr>
  {{end}}
</table>
```

`rowspan` offers a unique challenge. Namely, the spanned cell is only rendered
in the first row (the one with the rowspan attribute). From there on, that cell
is excluded in the HTML. Before looking at templating, this is an example of
correctly formatted HTML using rowspan:

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
the second row. That's because this field is spanned across two rows. This
presents a challenge for using a simple `range` command.

The solution is to use a couple of helpers from the
[sprig](http://masterminds.github.io/sprig/) library to render the first row
separately from the rest. Consider the following data that has been indexed by
model:

```json
{
  "devices_by_model": {
    "93180YC-EX": [
      { "name": "switch01", "model": "93180YC-EX" },
      { "name": "switch02", "model": "93180YC-EX" }
    ]
  }
}
```

```html
<table>
  <tr>
    <th>Model</th>
    <th>Switches</th>
  </tr>

  <!--render the first row, including the colspan-->
  {{first .devices_by_model}}
  <tr>
    <td rowspan="{{len .devices_by_model}}">{{.model}}</td>
    <td>{{.name}}</td>
  </tr>
  {{end}}

  <!--render the remaining rows-->
  {{range rest .devices_by_model}}
  <tr>
    <td>{{.name}}</td>
  </tr>
  {{end}}
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
automation for CX delivery. The following solutions were considered, and have
all been used in production leading up to this tool:

**SCDP** is our Atlasian Confluence deployment, and is a collaborative
templating, document creation, and document sharing platform. Document templates
can be reused and tailored to the needs of each project.

**Challenges:**

- SCDP is a document collaboration and templating solution; not an automation
  solution. There is no clear solution for automating _content_ into SCDP
  templates.

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
Markdown document (e.g. Jinja2, text/template, etc), and then Pandoc can be used
to transform this into a Microsoft Word document. This was the predecesor to
docgen and is the closest in architecture.

**Challenges:**

- Getting started, developing, and deploying a tool that uses this solution
  requires a lot more dependencies.
- Invalid Markdown or HTML can produce a broken/invalid Word document.
- Pandoc is not as fast or efficient as a purpose-built tool such as docgen.

### Feature Comparison

| Feature                                       | docgen             | SCDP               | Library            | Pandoc             |
| --------------------------------------------- | :----------------- | :----------------- | :----------------- | :----------------- |
| Automated                                     | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: |
| Separation of business and presentation logic | :white_check_mark: | **na**             | :x:                | :white_check_mark: |
| Declarative                                   | :white_check_mark: | **na**             | :x:                | :white_check_mark: |
| Easy to integrate with new tools              | :white_check_mark: | :x:                | :x: **¹**          | :x: **²**          |
| Git-friendly                                  | :white_check_mark: | **na**             | :x: **³**          | :white_check_mark: |
| Will not create a broken Word document        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Fast / memory efficient                       | :white_check_mark: | **na**             | :white_check_mark: | :x:                |
| Developer Friendly                            | :white_check_mark: | :x:                | :white_check_mark: | :x: **⁴**          |

**¹** A docx library such as python-docx can be very easy to get started with,
but writing _quality_ documentation in Python can be very challenging.

**²** Several independent components are required for this solution, making it
much more difficult for an end-user (conulting engineer) to run.

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

Integration tests validate that generated `.docx` files conform to the OOXML standard. See [INTEGRATION_TESTING.md](INTEGRATION_TESTING.md) for detailed instructions on running integration tests locally.

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

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
