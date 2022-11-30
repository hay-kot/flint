<br />
<div align="center">
  <h1 align="center">Flint</h1>
  <p align="center">
    Customizable FrontMatter linter for your static site generator
    <br />
    <a href="#configuration"><strong>Config Reference »</strong></a>
    <br />
    <br />
    <a href="https://github.com/hay-kot/flint/issues">Report Bug</a>
    ·
    <a href="https://github.com/hay-kot/flint/issues">Request Feature</a>
  </p>
</div>

Flint is a customizable linter for frontmatter that can can be tailored to your projects needs. It uses a configuration file to define rule sets that can be applied conditionally to files based on their path.

- Supported frontmatter formats:
  - YAML
  - JSON (Partial)
  - TOML (Partial)

| Format | Linting | Line / Col |
| ------ | :-----: | :--------: |
| YAML   |    ✓    |     ✓      |
| JSON   |    ✓    |     ✖      |
| TOML   |    ✓    |     ✖      |

This library is still experimental.

## Features

- Powerful glob matching support
- Customizable rule-sets
- Customizable error messages
- Powerful builtin checks
  - Required fields
  - Regex matching
  - Enums Sets
  - Date Format
  - Field Length
  - <s>Field types</s> (coming soon!)
  - <s>Asset Existence</s> (coming soon!)

## Installation

### Go Install

```bash
go install github.com/hay-kot/flint@latest
```

### Homebrew

```bash
brew tap hay-kot/flint-tap
brew install hay-kot/flint-tap/flint
```

## Usage

flint does nothing out of the box, you must initialize a configuration file in the root directory where you'll call flint. This is usually the root of your project.

### Configuration

The configuration has two areas (1) the `rules` and (2) the `content`. The rules is where you define a Key/Value pair of rules where the key is the ID of the rule and the value is the checks to be performed. The content is where you define the content to be linted and which rules to apply to it. You can write your configuration in YAML, JSON, or TOML, though we recommend YAML is it's the most readable.

```yaml
rules:
  FM001:
    level: error
    description: "required post fields are missing"
    required:
      - "title"
      - "date"

content:
  - name: Blog Posts
    paths:
      - content/**/*.{md,html}
    rules:
      - FM001
```

**Note:** If the configuration file cannot be found, flint will attempt to look for a glob match for a toml, json, or yaml file in the current directory, or the parent directory of the config path passed in. If you have multiple files that may match the glob, the first one found will be used. This may lead to unexpected results if you have multiple configuration files in your project. We recommend you use the `--config` flag to specify the path to your configuration file explicitly if you're using multiple configuration files.

#### Rules

The key is the ID of the rule. This is used to identify the rule in the output. We recommend using either `FM001` and up or define your own rule id using a common prefix and number scheme.

The value represents the configuration for the rule set being created. Flint comes with a set of built in rules that can be applied to keys within your front matter. Custom rules cannot be defined at this time. Note that if multiple checks are defined for a rule, they will all be applied to the front matter.

#### Properties

- `level` - The level of the rule. This can be `error`, `warning`, or `info`. This is used to determine the exit code of the program.
- `description` - A description of the rule. This is used in the output to describe the rule. Keep it short and sweet.

#### Built-in Rules

##### `bultin.required`

This rule checks to see if the keys defined in the array are present in the front matter. It's values is a list of strings

```yaml
rules:
  FM002:
    level: warning
    description: "author is missing"
    required:
      - "author.name"
      - "author.email"
```

##### `match`

This rules checks to see if the regular expressions defined in the array match the values of the fields defined in the array. You can provide multiple regular expressions and field, just keep in mind that every regex will be applied to each field.

Supported types are `string` and `[]string`

```yaml
  FM003:
    level: error
    description: "slug is not valid url"
    match:
      re:
        - "^[a-z0-9]+(?:-[a-z0-9]+)*$"
      fields:
        - "meta.slug"
        - "slug"
```

##### `enum`

checks to see if the string or string array are defined in the values.

Supported types are `string` and `[]string`

```yaml
  FM004:
    level: error
    description: "category not in list"
    enum:
      values:
        - "Go"
        - "Python"
      fields:
        - "category"
        - "post.category"
```

##### `date`

checks to see if the date is in the correct format. The format is defined using the [Go time package](https://golang.org/pkg/time/#pkg-constants)

Supported types are `string`

```yaml
  FM005:
    level: error
    description: "date is not in correct format"
    date:
      format:
        - "2006-01-02"
      fields:
        - "date"
```

##### `length`

Checks to see if the field meets a min/max length

Supported types are `string` and `[]string`

```yaml
  FM006:
    level: error
    description: "must be between 2 and 5"
    length:
      min: 2
      max: 5
      fields:
        - "keywords"
        - "categories"
        - "tags"
```

#### Content

The content is an array of objects that define the content to be linted and the rules to apply to it. You can use a robust glob pattern to define the paths to be linted. See [doublestar](https://github.com/bmatcuk/doublestar) for more details on supported patterns

```yaml
content:
  - name: Blog Posts
    paths:
      - content/**/*.{md,html}
    rules:
      - FM001
      - FM002
      - FM003
```

### Recipes

#### Blog Requirements

```yaml
rules:
  FM001:
    level: error
    description: "blog post requirements"
    required:
      - "title"
      - "slug"
      - "description"
      - "date"
      - "keywords"
    date:
      format:
        - "2006-01-02"
      fields:
        - "date"
    length:
      min: 2
      max: 5
      fields:
        - "keywords"
        - "categories"
        - "tags"
    match:
      re:
        - "^[a-z0-9]+(?:-[a-z0-9]+)*$"
      fields:
        - "slug"
```

#### Ignore Files that start with an underscore

```yaml
content:
  - name: Blog Posts
    paths:
      - content/blog/**/[!_]*.md
    rules:
      - FM001
      - FM002
      - FM003
```

#### Regular Expressions

| Description | Regex                      |
| ----------- | -------------------------- |
| Slug        | ^[a-z0-9]+(?:-[a-z0-9]+)*$ |
