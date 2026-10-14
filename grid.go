package layout

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Grid e um layout de grade que distribui filhos automaticamente em linhas,
// com colunas de largura igual.
type Grid struct {
	// Columns define o numero de colunas na grade.
	//
	// ATENÇÃO: o campo é público e pode ser alterado diretamente, o que
	// dispensa a validação feita por [Grid.NewGrid]. Use [Grid.Validate] (ou
	// [Grid.RenderE]) para checar a configuração atual e [Grid.ChildrenCopy]
	// para obter uma cópia que não afete a grade.
	Columns int
	// Children sao os conteudos dos filhos da grade.
	Children []string
	// Gap define o espacamento (em caracteres) entre celulas.
	Gap int

	width  int
	height int

	// err guarda o primeiro erro de configuração detectado pelos métodos de
	// construção.
	err error
}

// NewGrid cria uma nova grade com o numero de colunas especificado.
//
// Um número de colunas menor que 1 é registrado como erro de configuração e
// normalizado para 1; use [Grid.NewGridE] para receber o erro.
func NewGrid(columns int) *Grid {
	g := &Grid{
		Columns: columns,
		Gap:     0,
		width:   80,
		height:  24,
	}
	if columns < 1 {
		g.registrarErro(fmt.Errorf("%w: %d", ErrInvalidColumns, columns))
		g.Columns = 1
	}
	return g
}

// NewGridE cria uma nova grade e devolve erro quando o número de colunas é
// inválido.
func NewGridE(columns int) (*Grid, error) {
	g := NewGrid(columns)
	if err := g.Validate(); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSize define a largura e altura disponiveis para a grade.
//
// Valores negativos são registrados como erro de configuração; o valor
// informado é mantido e saneado na renderização.
func (g *Grid) SetSize(width, height int) *Grid {
	if width < 0 {
		g.registrarErro(fmt.Errorf("%w: largura %d", ErrInvalidSize, width))
	}
	if height < 0 {
		g.registrarErro(fmt.Errorf("%w: altura %d", ErrInvalidSize, height))
	}
	g.width = width
	g.height = height
	return g
}

// SetGap define o espacamento entre celulas da grade.
//
// Valores negativos são registrados como erro de configuração; o valor
// informado é mantido e saneado na renderização.
func (g *Grid) SetGap(gap int) *Grid {
	if gap < 0 {
		g.registrarErro(fmt.Errorf("%w: %d", ErrInvalidGap, gap))
	}
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

// Size devolve a largura e a altura disponíveis para a grade.
func (g *Grid) Size() (width, height int) {
	return g.width, g.height
}

// ChildCount devolve a quantidade de filhos da grade.
func (g *Grid) ChildCount() int {
	return len(g.Children)
}

// ChildrenCopy devolve uma cópia dos filhos da grade. Alterar a cópia não
// afeta a grade.
func (g *Grid) ChildrenCopy() []string {
	if len(g.Children) == 0 {
		return nil
	}
	copia := make([]string, len(g.Children))
	copy(copia, g.Children)
	return copia
}

// Validate checa a configuração atual da grade e devolve o primeiro erro
// encontrado.
func (g *Grid) Validate() error {
	if g.err != nil {
		return g.err
	}
	if g.Columns < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumns, g.Columns)
	}
	if g.Gap < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidGap, g.Gap)
	}
	if g.width < 0 {
		return fmt.Errorf("%w: largura %d", ErrInvalidSize, g.width)
	}
	if g.height < 0 {
		return fmt.Errorf("%w: altura %d", ErrInvalidSize, g.height)
	}
	return nil
}

// Render calcula o layout da grade e retorna como string.
// Distribui os filhos em linhas com colunas de largura igual.
//
// Esta variante é tolerante: nunca entra em panic e sempre devolve uma
// geometria válida, mesmo com configuração inválida. Use [Grid.RenderE] para
// detectar erros de configuração.
func (g *Grid) Render() string {
	return g.renderizar()
}

// RenderE calcula o layout e devolve erro quando a configuração da grade é
// inválida.
func (g *Grid) RenderE() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}
	return g.renderizar(), nil
}

// registrarErro guarda o primeiro erro de configuração encontrado.
func (g *Grid) registrarErro(err error) {
	if g.err == nil {
		g.err = err
	}
}

// renderizar monta a grade usando valores saneados.
func (g *Grid) renderizar() string {
	if len(g.Children) == 0 {
		return ""
	}

	colunas := max(g.Columns, 1)
	gap := max(g.Gap, 0)
	largura := max(g.width, 0)
	altura := max(g.height, 0)

	totalLinhas := (len(g.Children) + colunas - 1) / colunas

	// O espaçamento é reduzido quando não cabe nem na largura nem na altura,
	// garantindo que a grade nunca ultrapasse o container.
	gapEfetivo := gap
	if colunas > 1 {
		gapEfetivo = min(gapEfetivo, largura/(colunas-1))
	}
	if totalLinhas > 1 {
		gapEfetivo = min(gapEfetivo, altura/(totalLinhas-1))
	}

	disponivelLargura := largura - gapEfetivo*(colunas-1)
	disponivelAltura := altura - gapEfetivo*(totalLinhas-1)

	larguras := distribuir(disponivelLargura, pesosIguais(colunas))
	alturas := distribuir(disponivelAltura, pesosIguais(totalLinhas))
	for i := range alturas {
		// Cada linha ocupa no mínimo uma linha de terminal.
		alturas[i] = max(alturas[i], 1)
	}

	linhas := make([]string, 0, totalLinhas)

	for inicio := 0; inicio < len(g.Children); inicio += colunas {
		fim := min(inicio+colunas, len(g.Children))
		indiceLinha := inicio / colunas
		alturaLinha := alturas[indiceLinha]

		celulas := make([]string, 0, colunas)
		for j := inicio; j < fim; j++ {
			celulas = append(celulas, renderCelula(g.Children[j], larguras[j-inicio], alturaLinha))
		}

		// Preenche celulas vazias na ultima linha para manter alinhamento
		for j := len(celulas); j < colunas; j++ {
			celulas = append(celulas, renderCelula("", larguras[j], alturaLinha))
		}

		if gapEfetivo > 0 {
			linhas = append(linhas, joinWithSpacer(celulas, strings.Repeat(" ", gapEfetivo), lipgloss.Top))
			continue
		}
		linhas = append(linhas, lipgloss.JoinHorizontal(lipgloss.Top, celulas...))
	}

	if gapEfetivo > 0 {
		// O separador precisa de gap-1 quebras de linha: o bloco separador já
		// ocupa uma linha, o que faz o total de linhas em branco entre duas
		// linhas da grade ser exatamente igual ao gap.
		return joinVerticalWithSpacer(linhas, strings.Repeat("\n", gapEfetivo-1))
	}
	return lipgloss.JoinVertical(lipgloss.Left, linhas...)
}
