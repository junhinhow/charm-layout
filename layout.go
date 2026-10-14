// Package layout fornece um motor de layout flexbox para interfaces de terminal.
// Inspirado no CSS Flexbox, permite criar layouts responsivos com flex ratios,
// direcao (linha/coluna) e espacamento entre elementos.
//
// # Semântica das razões
//
// As razões dos filhos seguem a semântica de `flex-grow` do CSS: o espaço livre
// (a largura ou a altura do container menos os espaçamentos) é distribuído
// proporcionalmente às razões informadas. As razões não precisam somar 1:
// dois filhos com razão 0.3 cada dividem o espaço livre em 50/50.
//
// # Geometria exata
//
// O espaço disponível é sempre distribuído por inteiro, sem resto perdido: a
// soma das larguras (ou alturas) dos filhos mais os espaçamentos é exatamente a
// largura (ou altura) declarada em [Flex.SetSize] e [Grid.SetSize]. Conteúdo
// maior que a célula é quebrado pelo lipgloss e truncado quando não cabe.
//
// # Validação
//
// [Flex.Render] e [Grid.Render] são tolerantes: ajustam valores inválidos e
// devolvem sempre uma geometria válida, sem entrar em panic. [Flex.RenderE] e
// [Grid.RenderE] devolvem o erro de configuração encontrado.
package layout

import (
	"fmt"
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

// String devolve o nome da direção.
func (d Direction) String() string {
	switch d {
	case Row:
		return "Row"
	case Column:
		return "Column"
	default:
		return fmt.Sprintf("Direction(%d)", int(d))
	}
}

// Valida informa se a direção é conhecida pelo motor de layout.
func (d Direction) Valid() bool {
	return d == Row || d == Column
}

// FlexChild representa um filho dentro de um container flex,
// com seu conteudo e a proporcao de espaco que ocupa.
type FlexChild struct {
	// Ratio define a proporcao de espaco que este filho ocupa.
	// O valor funciona como `flex-grow` do CSS: o espaço livre do container é
	// distribuído proporcionalmente às razões. As razões não precisam somar 1
	// (dois filhos com razão 0.3 dividem o espaço em 50/50).
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
	//
	// ATENÇÃO: o campo é público e pode ser alterado diretamente, o que
	//dispensa a validação feita por [Flex.AddChild]. Use [Flex.Validate] (ou
	// [Flex.RenderE]) para checar a configuração atual e [Flex.ChildrenCopy]
	// para obter uma cópia que não afete o container.
	Children []FlexChild
	// Gap define o espacamento (em caracteres) entre filhos.
	Gap int

	width  int
	height int

	// err guarda o primeiro erro de configuração detectado pelos métodos de
	// construção, para que a validação continue disponível mesmo quando o
	// valor problemático já foi normalizado.
	err error
}

// NewFlex cria um novo container flex com a direcao especificada.
//
// Uma direção desconhecida é registrada como erro de configuração: use
// [Flex.Validate] ou [Flex.RenderE] para detectá-la. [Flex.Render] continua
// funcionando, tratando a direção inválida como [Column].
func NewFlex(dir Direction) *Flex {
	f := &Flex{
		Direction: dir,
		Gap:       0,
		width:     80,
		height:    24,
	}
	if !dir.Valid() {
		f.registrarErro(fmt.Errorf("%w: %d", ErrInvalidDirection, int(dir)))
	}
	return f
}

// NewFlexE cria um novo container flex e devolve erro quando a direção é
// desconhecida.
func NewFlexE(dir Direction) (*Flex, error) {
	f := NewFlex(dir)
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return f, nil
}

// SetSize define a largura e altura disponiveis para o layout.
// Deve ser chamado antes de Render() para layouts responsivos.
//
// Valores negativos são registrados como erro de configuração; o valor
// informado é mantido e saneado na renderização.
func (f *Flex) SetSize(width, height int) *Flex {
	if width < 0 {
		f.registrarErro(fmt.Errorf("%w: largura %d", ErrInvalidSize, width))
	}
	if height < 0 {
		f.registrarErro(fmt.Errorf("%w: altura %d", ErrInvalidSize, height))
	}
	f.width = width
	f.height = height
	return f
}

// SetGap define o espacamento entre filhos.
//
// Valores negativos são registrados como erro de configuração; o valor
// informado é mantido e saneado na renderização.
func (f *Flex) SetGap(gap int) *Flex {
	if gap < 0 {
		f.registrarErro(fmt.Errorf("%w: %d", ErrInvalidGap, gap))
	}
	f.Gap = gap
	return f
}

// AddChild adiciona um filho com o ratio e conteudo especificados.
//
// Uma razão negativa, NaN ou infinita é registrada como erro de
// configuração.
func (f *Flex) AddChild(ratio float64, content string) *Flex {
	if !ratioValido(ratio) {
		f.registrarErro(fmt.Errorf("%w: %v no filho %d", ErrInvalidRatio, ratio, len(f.Children)))
	}
	f.Children = append(f.Children, FlexChild{
		Ratio:   ratio,
		Content: content,
	})
	return f
}

// Size devolve a largura e a altura disponíveis para o layout.
func (f *Flex) Size() (width, height int) {
	return f.width, f.height
}

// ChildCount devolve a quantidade de filhos do container.
func (f *Flex) ChildCount() int {
	return len(f.Children)
}

// ChildrenCopy devolve uma cópia dos filhos do container. Alterar a cópia não
// afeta o container.
func (f *Flex) ChildrenCopy() []FlexChild {
	if len(f.Children) == 0 {
		return nil
	}
	copia := make([]FlexChild, len(f.Children))
	copy(copia, f.Children)
	return copia
}

// Validate checa a configuração atual do container e devolve o primeiro erro
// encontrado.
func (f *Flex) Validate() error {
	if f.err != nil {
		return f.err
	}
	if !f.Direction.Valid() {
		return fmt.Errorf("%w: %d", ErrInvalidDirection, int(f.Direction))
	}
	if f.Gap < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidGap, f.Gap)
	}
	if f.width < 0 {
		return fmt.Errorf("%w: largura %d", ErrInvalidSize, f.width)
	}
	if f.height < 0 {
		return fmt.Errorf("%w: altura %d", ErrInvalidSize, f.height)
	}

	soma := 0.0
	for i, child := range f.Children {
		if !ratioValido(child.Ratio) {
			return fmt.Errorf("%w: %v no filho %d", ErrInvalidRatio, child.Ratio, i)
		}
		soma += child.Ratio
	}
	if len(f.Children) > 0 && soma <= 0 {
		return ErrZeroRatioSum
	}

	return nil
}

// Render calcula os tamanhos de cada filho com base nos ratios e retorna
// o layout composto como uma string. Usa lipgloss.JoinHorizontal para Row
// e lipgloss.JoinVertical para Column.
//
// Esta variante é tolerante: nunca entra em panic e sempre devolve uma
// geometria válida, mesmo com configuração inválida. Use [Flex.RenderE] para
// detectar erros de configuração.
func (f *Flex) Render() string {
	saida, _ := f.renderizar()
	return saida
}

// RenderE calcula o layout e devolve erro quando a configuração do container
// é inválida.
func (f *Flex) RenderE() (string, error) {
	if err := f.Validate(); err != nil {
		return "", err
	}
	return f.renderizar()
}

// registrarErro guarda o primeiro erro de configuração encontrado.
func (f *Flex) registrarErro(err error) {
	if f.err == nil {
		f.err = err
	}
}

// renderizar monta o layout usando valores saneados.
func (f *Flex) renderizar() (string, error) {
	if len(f.Children) == 0 {
		return "", nil
	}

	direcao := f.Direction
	if !direcao.Valid() {
		// Fallback determinístico: direções desconhecidas são tratadas como
		// coluna, igual ao comportamento histórico do pacote.
		direcao = Column
	}

	gap := max(f.Gap, 0)
	largura := max(f.width, 0)
	altura := max(f.height, 0)

	pesos := make([]float64, len(f.Children))
	for i, child := range f.Children {
		pesos[i] = ratioSeguro(child.Ratio)
	}

	if direcao == Row {
		return f.renderizarLinha(gap, largura, altura, pesos)
	}
	return f.renderizarColuna(gap, largura, altura, pesos)
}

// renderizarLinha distribui os filhos horizontalmente.
func (f *Flex) renderizarLinha(gap, largura, altura int, pesos []float64) (string, error) {
	n := len(f.Children)

	// O espaçamento é reduzido quando não cabe na largura disponível.
	gapEfetivo := limitarGap(gap, largura, n)
	disponivel := largura - gapEfetivo*(n-1)

	larguras := distribuir(disponivel, pesos)
	alturaBloco := max(altura, 1)

	renderizados := make([]string, 0, n)
	for i, child := range f.Children {
		renderizados = append(renderizados, renderCelula(child.Content, larguras[i], alturaBloco))
	}

	if gapEfetivo > 0 {
		return joinWithSpacer(renderizados, strings.Repeat(" ", gapEfetivo), lipgloss.Top), nil
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderizados...), nil
}

// renderizarColuna distribui os filhos verticalmente.
func (f *Flex) renderizarColuna(gap, largura, altura int, pesos []float64) (string, error) {
	n := len(f.Children)

	// O espaçamento é reduzido quando não cabe na altura disponível.
	gapEfetivo := limitarGap(gap, altura, n)
	disponivel := altura - gapEfetivo*(n-1)

	alturas := distribuir(disponivel, pesos)
	for i := range alturas {
		// Cada filho ocupa no mínimo uma linha.
		alturas[i] = max(alturas[i], 1)
	}

	larguraBloco := max(largura, 0)

	renderizados := make([]string, 0, n)
	for i, child := range f.Children {
		renderizados = append(renderizados, renderCelula(child.Content, larguraBloco, alturas[i]))
	}

	if gapEfetivo > 0 {
		// O separador precisa de gap-1 quebras de linha: o bloco separador já
		// ocupa uma linha, o que faz o total de linhas em branco entre dois
		// filhos ser exatamente igual ao gap.
		return joinVerticalWithSpacer(renderizados, strings.Repeat("\n", gapEfetivo-1)), nil
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderizados...), nil
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
