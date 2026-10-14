# charm-layout

> [Leia em Português](README.pt-br.md)

A Flexbox-inspired layout engine for terminal UIs built on top of [lipgloss](https://charm.land/lipgloss/v2). Create responsive flex and grid layouts with proportional sizing, directional flow, and gap spacing — just like CSS Flexbox, but for the terminal.

## Features

- **Flex Layout** — Row or Column direction with proportional child sizing (ratios like 0.3 / 0.7)
- **Grid Layout** — Auto-distributes children into equal-width columns across rows
- **Gap Spacing** — Configurable spacing between children
- **Exact geometry** — The container is filled up to the exact declared width/height; no column is lost to rounding
- **Responsive** — `SetSize(w, h)` adapts the layout to any terminal dimensions
- **Validated** — Invalid configuration is reported (`Validate` / `RenderE`) instead of silently breaking the layout
- **Built on lipgloss** — Uses `JoinHorizontal` / `JoinVertical` internally for reliable rendering

## Installation

```bash
go get github.com/junhinhow/charm-layout
```

## Quick Start

### Flex Layout

```go
package main

import (
    "fmt"
    "github.com/junhinhow/charm-layout"
)

func main() {
    flex := layout.NewFlex(layout.Row).
        SetSize(80, 24).
        SetGap(1).
        AddChild(0.3, "Sidebar").
        AddChild(0.7, "Main Content")

    fmt.Println(flex.Render())
}
```

### Column Layout

```go
flex := layout.NewFlex(layout.Column).
    SetSize(80, 24).
    AddChild(0.2, "Header").
    AddChild(0.6, "Body").
    AddChild(0.2, "Footer")

fmt.Println(flex.Render())
```

### Grid Layout

```go
grid := layout.NewGrid(3).
    SetSize(90, 30).
    SetGap(1).
    AddChildren("Cell 1", "Cell 2", "Cell 3", "Cell 4", "Cell 5")

fmt.Println(grid.Render())
```

## Sizing semantics

Ratios follow the CSS `flex-grow` rule: the **free space** of the container
(its width — or height, for a column — minus the gaps) is distributed
proportionaly to the ratios.

- The ratios **do not need to add up to 1**. Two children with ratio `0.3`
  split the free space 50/50.
- The free space is distributed with the *largest remainder* method, so the
  children always fill the container **exactly**: no column is lost to
  rounding and no space is left over.
- A ratio of `0` gives the child zero width; it is still rendered (and
  truncated to zero columns) so the layout stays aligned.
- If every ratio is zero, the space is split equally between the children.
- Every child (or grid row) occupies at least one terminal line.

```go
// 80 columns, gap 1, two children of ratio 0.3 each:
// the first child occupies columns 0..39, the second 40..79.
flex := layout.NewFlex(layout.Row).SetSize(80, 24).SetGap(1).
    AddChild(0.3, "A").AddChild(0.3, "B")
```

### Grid geometry

The columns share the available width equally, and the remainder of the
division is given to the first columns, so the row always ends exactly at the
declared width:

```go
// 80 columns, 10 columns, gap 1: the columns are 7 or 8 cells wide and the
// row ends exactly at column 80.
grid := layout.NewGrid(10).SetSize(80, 24).SetGap(1).
    AddChildren("A", "B", "C", "D", "E", "F", "G", "H", "I", "J")
```

When the gap does not fit in the available space it is reduced to the largest
value that still fits, so the layout never overflows the terminal.

## Validation

`Render()` is **tolerant**: it never panics and always returns a valid
geometry, sanitizing invalid values (negative sizes, negative gaps, unknown
directions, invalid ratios).

`RenderE()` is **strict**: it returns the configuration error and renders
nothing. Use it when the layout is built from user input.

```go
flex := layout.NewFlex(layout.Row).SetSize(80, 24).AddChild(-0.5, "A")

if _, err := flex.RenderE(); err != nil {
    // err matches errors.Is(err, layout.ErrInvalidRatio)
}

err := flex.Validate() // same check, without rendering
```

The following sentinel errors are returned (always wrapped with context, so
`errors.Is` works):

| Error | Meaning |
|-------|---------|
| `ErrInvalidDirection` | Unknown flex direction |
| `ErrInvalidGap` | Negative gap |
| `ErrInvalidSize` | Negative width or height |
| `ErrInvalidRatio` | Negative, `NaN` or infinite ratio |
| `ErrZeroRatioSum` | All ratios add up to zero |
| `ErrInvalidColumns` | Grid with less than one column |

## API Reference

### Flex

| Method | Description |
|--------|-------------|
| `NewFlex(dir Direction) *Flex` | Create a new flex container (Row or Column) |
| `NewFlexE(dir Direction) (*Flex, error)` | Same, but fails on an unknown direction |
| `SetSize(width, height int) *Flex` | Set available dimensions |
| `SetGap(gap int) *Flex` | Set spacing between children |
| `AddChild(ratio float64, content string) *Flex` | Add a child with proportional sizing |
| `Render() string` | Compute layout and return composed string (tolerant) |
| `RenderE() (string, error)` | Compute layout, returning the configuration error |
| `Validate() error` | Check the current configuration |
| `ChildrenCopy() []FlexChild` | Copy of the children slice |
| `Size() (int, int)` | Available width and height |
| `ChildCount() int` | Number of children |

### Grid

| Method | Description |
|--------|-------------|
| `NewGrid(columns int) *Grid` | Create a grid with N columns |
| `NewGridE(columns int) (*Grid, error)` | Same, but fails on an invalid column count |
| `SetSize(width, height int) *Grid` | Set available dimensions |
| `SetGap(gap int) *Grid` | Set spacing between cells |
| `AddChild(content string) *Grid` | Add a single child |
| `AddChildren(contents ...string) *Grid` | Add multiple children |
| `Render() string` | Compute grid layout and return composed string (tolerant) |
| `RenderE() (string, error)` | Compute grid layout, returning the configuration error |
| `Validate() error` | Check the current configuration |
| `ChildrenCopy() []string` | Copy of the children slice |
| `Size() (int, int)` | Available width and height |
| `ChildCount() int` | Number of children |

### Types

```go
type Direction int

const (
    Row    Direction = iota  // Horizontal distribution
    Column                   // Vertical distribution
)

func (d Direction) String() string  // "Row", "Column" or "Direction(n)"
func (d Direction) Valid() bool    // true for Row and Column

type FlexChild struct {
    Ratio   float64  // flex-grow weight — the ratios do not need to add up to 1
    Content string   // Text content
}
```

## Comparison with CSS Flexbox

| CSS Flexbox | charm-layout |
|-------------|-------------|
| `display: flex` | `NewFlex(Row)` |
| `flex-direction: column` | `NewFlex(Column)` |
| `flex: 0.3` | `AddChild(0.3, content)` |
| `gap: 8px` | `SetGap(1)` |
| `display: grid; grid-template-columns: repeat(3, 1fr)` | `NewGrid(3)` |

## License

MIT
