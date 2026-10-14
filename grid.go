package layout

import (
	"charm.land/lipgloss/v2"
)

// Grid e um layout de grade que distribui filhos automaticamente em linhas,
// com colunas de largura igual.
type Grid struct {
	// Columns define o numero de colunas na grade.
	Columns int
	// Children sao os conteudos dos filhos da grade.
	Children []string
	// Gap define o espacamento (em caracteres) entre celulas.
	Gap int

	width  int
	height int
}

// NewGrid cria uma nova grade com o numero de colunas especificado.
func NewGrid(columns int) *Grid {
	if columns < 1 {
		columns = 1
	}
	return &Grid{
		Columns: columns,
		Gap:     0,
		width:   80,
		height:  24,
	}
}

// SetSize define a largura e altura disponiveis para a grade.
func (g *Grid) SetSize(width, height int) *Grid {
	g.width = width
	g.height = height
	return g
}

// SetGap define o espacamento entre celulas da grade.
func (g *Grid) SetGap(gap int) *Grid {
	g.Gap = gap
	return g
}

// AddChild adiciona um filho a grade.
func (g *Grid) AddChild(content string) *Grid {
	g.Children = append(g.Children, content)
	return g
}

// AddChildren adiciona multiplos filhos a grade de uma vez.
func (g *Grid) AddChildren(contents ...string) *Grid {
	g.Children = append(g.Children, contents...)
	return g
}

// Render calcula o layout da grade e retorna como string.
// Distribui os filhos em linhas com colunas de largura igual.
func (g *Grid) Render() string {
	if len(g.Children) == 0 {
		return ""
	}

	// Columns e um campo publico: pode ser alterado depois da criacao.
	// Valida aqui para evitar divisao por zero e parametros invalidos.
	columns := g.Columns
	if columns < 1 {
		columns = 1
	}

	gap := g.Gap
	if gap < 0 {
		gap = 0
	}

	width := g.width
	if width < 0 {
		width = 0
	}

	height := g.height
	if height < 0 {
		height = 0
	}

	totalHGap := gap * (columns - 1)
	colWidth := (width - totalHGap) / columns
	if colWidth < 0 {
		colWidth = 0
	}

	// Calcula numero de linhas
	rowCount := (len(g.Children) + columns - 1) / columns
	totalVGap := gap * (rowCount - 1)
	rowHeight := 0
	if rowCount > 0 {
		rowHeight = (height - totalVGap) / rowCount
	}
	if rowHeight < 0 {
		rowHeight = 0
	}

	cellStyle := lipgloss.NewStyle().
		Width(colWidth).
		Height(rowHeight)

	rows := make([]string, 0, rowCount)

	for i := 0; i < len(g.Children); i += columns {
		end := i + columns
		if end > len(g.Children) {
			end = len(g.Children)
		}

		cells := make([]string, 0, columns)
		for j := i; j < end; j++ {
			cells = append(cells, cellStyle.Render(g.Children[j]))
		}

		// Preenche celulas vazias na ultima linha para manter alinhamento
		for len(cells) < columns {
			cells = append(cells, cellStyle.Render(""))
		}

		if gap > 0 {
			row := joinWithSpacer(cells, spacerH(gap), lipgloss.Top)
			rows = append(rows, row)
		} else {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		}
	}

	if gap > 0 {
		return joinVerticalWithSpacer(rows, spacerV(gap))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// spacerH retorna um espacador horizontal com a largura especificada.
func spacerH(width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s += " "
	}
	return s
}

// spacerV retorna um espacador vertical com a altura especificada.
func spacerV(height int) string {
	s := ""
	for i := 0; i < height; i++ {
		s += "\n"
	}
	return s
}
