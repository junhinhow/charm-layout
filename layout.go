// Package layout fornece um motor de layout flexbox para interfaces de terminal.
// Inspirado no CSS Flexbox, permite criar layouts responsivos com flex ratios,
// direcao (linha/coluna) e espacamento entre elementos.
package layout

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Direction define a direcao do layout flex.
type Direction int

const (
	// Row distribui os filhos horizontalmente.
	Row Direction = iota
	// Column distribui os filhos verticalmente.
	Column
)

// FlexChild representa um filho dentro de um container flex,
// com seu conteudo e a proporcao de espaco que ocupa.
type FlexChild struct {
	// Ratio define a proporcao de espaco que este filho ocupa (0.0 a 1.0).
	// Exemplo: 0.3 = 30% do espaco disponivel.
	Ratio float64
	// Content e o conteudo textual do filho.
	Content string
}

// Flex e um container de layout flexivel que distribui seus filhos
// de acordo com ratios proporcionais, similar ao CSS Flexbox.
type Flex struct {
	// Direction define se os filhos sao distribuidos em linha ou coluna.
	Direction Direction
	// Children sao os filhos do container flex.
	Children []FlexChild
	// Gap define o espacamento (em caracteres) entre filhos.
	Gap int

	width  int
	height int
}

// NewFlex cria um novo container flex com a direcao especificada.
func NewFlex(dir Direction) *Flex {
	return &Flex{
		Direction: dir,
		Gap:       0,
		width:     80,
		height:    24,
	}
}

// SetSize define a largura e altura disponiveis para o layout.
// Deve ser chamado antes de Render() para layouts responsivos.
func (f *Flex) SetSize(width, height int) *Flex {
	f.width = width
	f.height = height
	return f
}

// SetGap define o espacamento entre filhos.
func (f *Flex) SetGap(gap int) *Flex {
	f.Gap = gap
	return f
}

// AddChild adiciona um filho com o ratio e conteudo especificados.
func (f *Flex) AddChild(ratio float64, content string) *Flex {
	f.Children = append(f.Children, FlexChild{
		Ratio:   ratio,
		Content: content,
	})
	return f
}

// Render calcula os tamanhos de cada filho com base nos ratios e retorna
// o layout composto como uma string. Usa lipgloss.JoinHorizontal para Row
// e lipgloss.JoinVertical para Column.
func (f *Flex) Render() string {
	if len(f.Children) == 0 {
		return ""
	}

	totalGap := f.Gap * (len(f.Children) - 1)
	rendered := make([]string, 0, len(f.Children))

	for _, child := range f.Children {
		var style lipgloss.Style

		if f.Direction == Row {
			availableWidth := f.width - totalGap
			childWidth := int(float64(availableWidth) * child.Ratio)
			style = lipgloss.NewStyle().
				Width(childWidth).
				Height(f.height)
		} else {
			availableHeight := f.height - totalGap
			childHeight := int(float64(availableHeight) * child.Ratio)
			style = lipgloss.NewStyle().
				Width(f.width).
				Height(childHeight)
		}

		rendered = append(rendered, style.Render(child.Content))
	}

	if f.Direction == Row {
		if f.Gap > 0 {
			spacer := strings.Repeat(" ", f.Gap)
			return joinWithSpacer(rendered, spacer, lipgloss.Top)
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	}

	if f.Gap > 0 {
		spacer := strings.Repeat("\n", f.Gap)
		return joinVerticalWithSpacer(rendered, spacer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

// joinWithSpacer une strings horizontalmente inserindo um espacador entre elas.
func joinWithSpacer(items []string, spacer string, pos lipgloss.Position) string {
	if len(items) == 0 {
		return ""
	}

	interleaved := make([]string, 0, len(items)*2-1)
	for i, item := range items {
		if i > 0 {
			interleaved = append(interleaved, spacer)
		}
		interleaved = append(interleaved, item)
	}
	return lipgloss.JoinHorizontal(pos, interleaved...)
}

// joinVerticalWithSpacer une strings verticalmente inserindo um espacador entre elas.
func joinVerticalWithSpacer(items []string, spacer string) string {
	if len(items) == 0 {
		return ""
	}

	interleaved := make([]string, 0, len(items)*2-1)
	for i, item := range items {
		if i > 0 {
			interleaved = append(interleaved, spacer)
		}
		interleaved = append(interleaved, item)
	}
	return lipgloss.JoinVertical(lipgloss.Left, interleaved...)
}
