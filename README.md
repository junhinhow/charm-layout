# charm-layout

> [Leia em Portugues](README.pt-br.md)

A Flexbox-inspired layout engine for terminal UIs built on top of [lipgloss](https://charm.land/lipgloss/v2). Create responsive flex and grid layouts with proportional sizing, directional flow, and gap spacing — just like CSS Flexbox, but for the terminal.

## Features

- **Flex Layout** — Row or Column direction with proportional child sizing (ratios like 0.3 / 0.7)
- **Grid Layout** — Auto-distributes children into equal-width columns across rows
- **Gap Spacing** — Configurable spacing between children
- **Responsive** — `SetSize(w, h)` adapts the layout to any terminal dimensions
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

## API Reference

### Flex

| Method | Description |
|--------|-------------|
| `NewFlex(dir Direction) *Flex` | Create a new flex container (Row or Column) |
| `SetSize(width, height int) *Flex` | Set available dimensions |
| `SetGap(gap int) *Flex` | Set spacing between children |
| `AddChild(ratio float64, content string) *Flex` | Add a child with proportional sizing |
| `Render() string` | Compute layout and return composed string |

### Grid

| Method | Description |
|--------|-------------|
| `NewGrid(columns int) *Grid` | Create a grid with N columns |
| `SetSize(width, height int) *Grid` | Set available dimensions |
| `SetGap(gap int) *Grid` | Set spacing between cells |
| `AddChild(content string) *Grid` | Add a single child |
| `AddChildren(contents ...string) *Grid` | Add multiple children |
| `Render() string` | Compute grid layout and return composed string |

### Types

```go
type Direction int

const (
    Row    Direction = iota  // Horizontal distribution
    Column                   // Vertical distribution
)

type FlexChild struct {
    Ratio   float64  // 0.0 to 1.0 — proportion of available space
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
