# Flint - Extensible Linter for Front Matter

Flint is a customizable linter for frontmatter that can can be tailed to your projects needs. It uses a configuration file to define rule sets that can be applied conditionally to files based on their path.

- Supported frontmatter formats:
  - YAML
  - <s>JSON</s> (coming soon!)
  - <s>TOML</s> (coming soon!)

This library is still experimental and probably broken. Use at your own risk.

## Features

- Powerful glob matching support
- Customizable rule-sets
- Customizable error messages
- Powerful builtin checks
  - Required fields
  - Regex matching
  - Enums Sets
  - <s>Field types</s> (coming soon!)
  - <s>Field length</s> (coming soon!)
  - <s>Parsable Date</s> (coming soon!)
  - <s>Date Format</s> (coming soon!)
  - <s>Asset Existence</s> (coming soon!)

## Installation

// TODO

## Usage

flint does nothing out of the box, you must initialize a configuration file in the root directory where you'll call flint. This is usually the root of your project.

### Configuration

The configuration has two areas (1) the `rules` and (2) the `content`. The rules is where you define a Key/Value pair of rules where the key is the ID of the rule and the value is the checks to be performed. The content is where you define the content to be linted and which rules to apply to it.

```yaml
rules:
  FM001:
    level: error
    description: "required post fields are missing"
    builtin.required:
      - "title"
      - "date"

content:
  - name: Blog Posts
    paths:
      - content/**/*.{md,html}
    rules:
      - FM001
```

#### Rules

##### Key

The key is the ID of the rule. This is used to identify the rule in the output. We recommend using either `FM001` and up or define your own rule id using a common prefix and number scheme.

##### Value

The value represents the configuration for the rule set being created. Flint comes with a set of built in rules that can be applied to keys within your front matter. Custom rules cannot be defined at this time. Note that if multiple checks are defined for a rule, they will all be applied to the front matter.

###### Properties

- `level` - The level of the rule. This can be `error`, `warning`, or `info`. This is used to determine the exit code of the program.
- `description` - A description of the rule. This is used in the output to describe the rule. Keep it short and sweet.

`bultin.required` - This rule checks to see if the keys defined in the array are present in the front matter. It's values is a list of strings

```yaml
rules:
  FM002:
    level: warning
    description: "author is missing"
    builtin.required:
      - "author.name"
      - "author.email"
```

`builtin.match` - This rules checks to see if the regular expressions defined in the array match the values of the fields defined in the array. You can provide mutliple regular expressions and field, just keep in mind that every regex will be applied to each field.

```yaml
  FM003:
    level: error
    description: "slug is not valid url"
    builtin.match:
      re:
        - "^[a-z0-9]+(?:-[a-z0-9]+)*$"
      fields:
        - "meta.slug"
        - "slug"
```

`builtin.enum` - checks to see if the string or string array are defined in the values.

Supported types are `string` and `[]string`

```yaml
  FM004:
    level: error
    description: "category not in list"
    builtin.enum:
      values:
        - "Go"
        - "Python"
      fields:
        - "category"
        - "post.category"
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

TODO
