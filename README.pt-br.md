# charm-layout

> [Read in English](README.md)

Um motor de layout inspirado no Flexbox para interfaces de terminal, construido sobre o [lipgloss](https://charm.land/lipgloss/v2). Crie layouts flex e grid responsivos com dimensionamento proporcional, direcao de fluxo e espacamento — assim como o CSS Flexbox, mas para o terminal.

## Funcionalidades

- **Layout Flex** — Direcao em Linha ou Coluna com dimensionamento proporcional dos filhos (ratios como 0.3 / 0.7)
- **Layout Grid** — Distribui filhos automaticamente em colunas de largura igual
- **Espacamento (Gap)** — Espacamento configuravel entre filhos
- **Responsivo** — `SetSize(w, h)` adapta o layout para qualquer dimensao de terminal
- **Construido sobre lipgloss** — Usa `JoinHorizontal` / `JoinVertical` internamente para renderizacao confiavel

## Instalacao

```bash
go get github.com/junhinhow/charm-layout
```

## Inicio Rapido

### Layout Flex

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
        AddChild(0.3, "Barra Lateral").
        AddChild(0.7, "Conteudo Principal")

    fmt.Println(flex.Render())
}
```

### Layout em Coluna

```go
flex := layout.NewFlex(layout.Column).
    SetSize(80, 24).
    AddChild(0.2, "Cabecalho").
    AddChild(0.6, "Corpo").
    AddChild(0.2, "Rodape")

fmt.Println(flex.Render())
```

### Layout Grid

```go
grid := layout.NewGrid(3).
    SetSize(90, 30).
    SetGap(1).
    AddChildren("Celula 1", "Celula 2", "Celula 3", "Celula 4", "Celula 5")

fmt.Println(grid.Render())
```

## Referencia da API

### Flex

| Metodo | Descricao |
|--------|-----------|
| `NewFlex(dir Direction) *Flex` | Cria um container flex (Row ou Column) |
| `SetSize(width, height int) *Flex` | Define as dimensoes disponiveis |
| `SetGap(gap int) *Flex` | Define o espacamento entre filhos |
| `AddChild(ratio float64, content string) *Flex` | Adiciona um filho com dimensionamento proporcional |
| `Render() string` | Calcula o layout e retorna a string composta |

### Grid

| Metodo | Descricao |
|--------|-----------|
| `NewGrid(columns int) *Grid` | Cria uma grade com N colunas |
| `SetSize(width, height int) *Grid` | Define as dimensoes disponiveis |
| `SetGap(gap int) *Grid` | Define o espacamento entre celulas |
| `AddChild(content string) *Grid` | Adiciona um unico filho |
| `AddChildren(contents ...string) *Grid` | Adiciona multiplos filhos |
| `Render() string` | Calcula o layout da grade e retorna a string composta |

### Tipos

```go
type Direction int

const (
    Row    Direction = iota  // Distribuicao horizontal
    Column                   // Distribuicao vertical
)

type FlexChild struct {
    Ratio   float64  // 0.0 a 1.0 — proporcao do espaco disponivel
    Content string   // Conteudo textual
}
```

## Comparacao com CSS Flexbox

| CSS Flexbox | charm-layout |
|-------------|-------------|
| `display: flex` | `NewFlex(Row)` |
| `flex-direction: column` | `NewFlex(Column)` |
| `flex: 0.3` | `AddChild(0.3, conteudo)` |
| `gap: 8px` | `SetGap(1)` |
| `display: grid; grid-template-columns: repeat(3, 1fr)` | `NewGrid(3)` |

## Licenca

MIT
