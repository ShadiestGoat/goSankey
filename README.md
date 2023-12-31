# Go Sankey

This is a tool that can generate a sankey chart from a single file.

## Why?

This was initially created as a 1am idea for a sleepless night, until it turned into a 24h project challenge. I plan to provide support for bug fixes/expansions on certain features.

## Usage

This software fully relies on a config file & directory. The easiest way to get started is by using the CLI:

1. `go install github.com/shadiestgoat/goSankey`
2. Navigate into a new empty folder
3. `goSankey init` - This will create a folder called "resources" and a config file "config.sankey"
4. Edit the config file (see later docs)
5. `goSankey create config.sankey` (or replace `config.sankey`, with the actual path to your config file)

## Config File

As you may have noticed, the config file is a bit special. This was an intentional part of the challenge. The config file follows the general principals of a `.conf` or `.ini` file, with some exceptions. Here is an overview of the syntax:

- Empty lines are ignored
- Comments can start with a `#` or with a `;`
- *Sections* must be surrounded by `[[]]`
- In the `config` and the `nodes` section, key-value pairs are represented as `key=value`
  - In here, both the key and the value are whitespace trimmed
  - Also, the key is case insensitive
  - If 2 of the same key are present, the later one is used 
  - Keys are sometimes documented with *PascalCase*, these keys are also available with *snake_case* and *with spaces*
  - Values can take up multiple *types*:
    - `INT` - Integer
    - `COLOR` - A hexadecimal color.
      - This supports 3 letter values like `fff` -> `ffffff`
      - There is semantic support for 8 letter values, but alpha channels are ignored `FFFFFF22` -> `FFFFFF`
      - Prefixes like `#` or `0x` are supported
      - Case insensitive
      - Sometimes this can be set to random, which creates a random contrasting color to the background
    - `BOOL` - A boolean
      - Case insensitive
      - `Yes`, `True`, `1` - true
      - `No`, `False`, `0` - false
    - `TEXT` - A string or text
- Sections can be broken up by other sections. The parser will 'glue' them back together
- Section names are case insensitive
- Each section has some special syntax rules & setup, see next section

### Section - Config
The config section has general configuration for your setup. This is a pure key-value store:

```
[[Config]]
KEY=Value
otHer_Key   = YES
```

The following values are configured:

| Key | Type | Description | Default |
|:---:|:----:|-------------|---------|
| ConnectionOpacity | INT | A percentage value that indicates the opacity of the node connections | 20 |
| Background | COLOR | The color of the background | #F6F8FA |
| Width | INT | The width, in pixels, of the output image | If height is present, then a value that formats the image as a `16:9` ratio image, otherwise, 1920 |
| Height | INT | The height, in pixels, of the output image | If width is present, then a value that formats the image as a `16:9` ratio image, otherwise, 1080 |
| OutputName | TEXT | The path to the output image. If empty, then a 'dry run' is performed. Note: no matter the extension, the app will always write a PNG into it! | Empty |
| DrawBorder | BOOL | If true, it will draw a border around the entire image (see BorderSize and BorderColor) | Yes |
| BorderColor | COLOR | What color should the border be? | Random |
| BorderSize | INT | The size, in pixels, of the border (note: this does not affect any positioning) | 2 |
| BorderPadding | INT | A percentage value showing how much padding there should be around the entire image (affects positioning, has nothing to do with DrawBorder) | 2 |
| NodeWidth | INT | A percentage value indicating the width of the width of the nodes | 2 |
| PadLeft | INT | A percentage value indicating how much extra padding there should be on the left (used in case the left tiles are just so big) | 1 |
| VertSpaceNodes | INT | A percentage value indicating how much height in total the nodes should take | 85 |
| HorizontalTextPad | INT | The amount, in pixels, that text boxes should be padded with horizontally | 15 |
| TextLinePad | INT | The amount, in pixels, that each line of a text should pad below itself | 5 |

### Section - Nodes
This section describes nodes, (ie. the bars). This section's special mechanic is that each node is described with an ID, surrounded by `[]`:

```
[[Nodes]]
[ID 1]
Key = Value
Key 2 = Value 2

[ID 2]
Key = Value
Key 2 = Value
```

Each ID is a sub-section of it's own, which means it doesn't share any keys with other sub-sections. Be careful - this ID is case-sensitive! It will be used in the next section. Each node is described by it's keys & values:

| Key | Type | Description | Default |
|:---:|:----:|-------------|---------|
| Title | TEXT | The human readable title of the node | The value of the ID |
| Color | COLOR | By default, it creates a random (but contrasting to the background) color |
| Step | INT | This indicates a horizontal position for the node. The steps are normalized by the parser, so don't worry about creating consecutive values | **This is a required value** |

### Section - Connections
This section is the most special. Each line follows the following format:
`{ORIGIN ID} -> {DESTINATION ID}: {AMOUNT}`
This section indicates the connections between nodes & their amount.

## Changing font

Don't like a font I'm using? Change it. Simply create a folder called "resources" in the working directory. In that directory, place the fonts you want to be used (alphabetically the first valid file will be the main font, with the rest being fallbacks)

## Example

![Example Output](https://github.com/ShadiestGoat/goSankey/blob/main/example.png?raw=true)