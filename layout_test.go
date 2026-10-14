package layout

import (
	"errors"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// colunaDe devolve a coluna (largura de exibição) em que texto aparece na
// linha informada.
func colunaDe(linha, texto string) int {
	idx := strings.Index(linha, texto)
	if idx < 0 {
		return -1
	}
	return ansi.StringWidth(linha[:idx])
}

// primeiraLinha devolve a primeira linha da saída renderizada.
func primeiraLinha(saida string) string {
	linhas := strings.Split(saida, "\n")
	return linhas[0]
}

func TestFlexLarguraExata(t *testing.T) {
	flex := NewFlex(Row).
		SetSize(40, 3).
		AddChild(0.3, "A").
		AddChild(0.7, "B")

	saida := flex.Render()
	if got := lipgloss.Width(saida); got != 40 {
		t.Errorf("largura renderizada = %d, esperado 40", got)
	}
}

func TestFlexRazaoNaoPrecisaSomarUm(t *testing.T) {
	// Semântica de flex-grow: o espaço livre é distribuído proporcionalmente
	// às razões, que não precisam somar 1.
	flex := NewFlex(Row).
		SetSize(80, 3).
		AddChild(0.3, "A").
		AddChild(0.3, "B")

	linha := primeiraLinha(flex.Render())

	if got := colunaDe(linha, "A"); got != 0 {
		t.Errorf("coluna do primeiro filho = %d, esperado 0", got)
	}
	if got := colunaDe(linha, "B"); got != 40 {
		t.Errorf("coluna do segundo filho = %d, esperado 40 (divisão 50/50)", got)
	}
	if got := lipgloss.Width(linha); got != 80 {
		t.Errorf("largura renderizada = %d, esperado 80", got)
	}
}

func TestFlexRazaoInvalida(t *testing.T) {
	flex := NewFlex(Row).SetSize(20, 2).AddChild(-0.5, "A")

	if _, err := flex.RenderE(); !errors.Is(err, ErrRatioInvalido) {
		t.Fatalf("RenderE devolveu %v, esperado ErrRatioInvalido", err)
	}

	// A variante tolerante não entra em panic e devolve uma geometria válida.
	saida := flex.Render()
	if got := lipgloss.Width(saida); got != 20 {
		t.Errorf("largura renderizada = %d, esperado 20", got)
	}
}

func TestFlexDirecaoInvalida(t *testing.T) {
	flex := NewFlex(Direction(42)).SetSize(20, 2).AddChild(0.5, "A")

	if _, err := flex.RenderE(); !errors.Is(err, ErrDirecaoInvalida) {
		t.Fatalf("RenderE devolveu %v, esperado ErrDirecaoInvalida", err)
	}
	if err := flex.Validate(); !errors.Is(err, ErrDirecaoInvalida) {
		t.Fatalf("Validate devolveu %v, esperado ErrDirecaoInvalida", err)
	}

	// A variante tolerante não entra em panic.
	if flex.Render() == "" {
		t.Error("Render deveria devolver uma saída não vazia")
	}
}

func TestFlexNewFlexE(t *testing.T) {
	if _, err := NewFlexE(Direction(7)); !errors.Is(err, ErrDirecaoInvalida) {
		t.Fatalf("NewFlexE devolveu %v, esperado ErrDirecaoInvalida", err)
	}
	if _, err := NewFlexE(Row); err != nil {
		t.Fatalf("NewFlexE(Row) devolveu %v, esperado nil", err)
	}
}

func TestFlexGapRespeitado(t *testing.T) {
	flex := NewFlex(Row).
		SetSize(30, 2).
		SetGap(3).
		AddChild(0.5, "A").
		AddChild(0.5, "B")

	linha := primeiraLinha(flex.Render())

	if got := colunaDe(linha, "B"); got != 17 {
		t.Errorf("coluna do segundo filho = %d, esperado 17", got)
	}
	if got := lipgloss.Width(linha); got != 30 {
		t.Errorf("largura renderizada = %d, esperado 30", got)
	}
}

func TestFlexGapMaiorQueEspaco(t *testing.T) {
	// O gap é reduzido para caber na largura: o layout não pode estourar.
	flex := NewFlex(Row).
		SetSize(10, 2).
		SetGap(50).
		AddChild(0.5, "A").
		AddChild(0.5, "B")

	saida := flex.Render()
	if got := lipgloss.Width(saida); got > 10 {
		t.Errorf("largura renderizada = %d, esperado no máximo 10", got)
	}
}

func TestFlexColunaAlturaExata(t *testing.T) {
	flex := NewFlex(Column).
		SetSize(20, 10).
		SetGap(2).
		AddChild(0.5, "A").
		AddChild(0.5, "B")

	saida := flex.Render()
	if got := lipgloss.Height(saida); got != 10 {
		t.Errorf("altura renderizada = %d, esperado 10", got)
	}
	if got := lipgloss.Width(saida); got != 20 {
		t.Errorf("largura renderizada = %d, esperado 20", got)
	}
}

func TestFlexConteudoLongoNaoEstoura(t *testing.T) {
	flex := NewFlex(Row).
		SetSize(12, 2).
		AddChild(0.5, "conteúdo bem maior que a célula")

	saida := flex.Render()
	if got := lipgloss.Width(saida); got != 12 {
		t.Errorf("largura renderizada = %d, esperado 12", got)
	}
}

func TestFlexConteudoMultibyteNaoEstoura(t *testing.T) {
	flex := NewFlex(Row).
		SetSize(11, 2).
		AddChild(0.5, "日本語日本語")

	saida := flex.Render()
	if got := lipgloss.Width(saida); got != 11 {
		t.Errorf("largura renderizada = %d, esperado 11", got)
	}
}

func TestFlexSetSizeNegativo(t *testing.T) {
	flex := NewFlex(Row).SetSize(-10, 5).AddChild(0.5, "A")
	if _, err := flex.RenderE(); !errors.Is(err, ErrTamanhoInvalido) {
		t.Fatalf("RenderE devolveu %v, esperado ErrTamanhoInvalido", err)
	}
}

func TestFlexSetGapNegativo(t *testing.T) {
	flex := NewFlex(Row).SetGap(-1).AddChild(0.5, "A")
	if _, err := flex.RenderE(); !errors.Is(err, ErrGapInvalido) {
		t.Fatalf("RenderE devolveu %v, esperado ErrGapInvalido", err)
	}
	if got := lipgloss.Width(flex.Render()); got != 80 {
		t.Errorf("largura renderizada = %d, esperado 80", got)
	}
}

func TestFlexSomaRazaoZero(t *testing.T) {
	flex := NewFlex(Row).SetSize(20, 2).AddChild(0, "A").AddChild(0, "B")
	if _, err := flex.RenderE(); !errors.Is(err, ErrRazaoSomaZero) {
		t.Fatalf("RenderE devolveu %v, esperado ErrRazaoSomaZero", err)
	}
	// A variante tolerante divide o espaço igualmente.
	linha := primeiraLinha(flex.Render())
	if got := colunaDe(linha, "B"); got != 10 {
		t.Errorf("coluna do segundo filho = %d, esperado 10", got)
	}
}

func TestFlexChildrenCopy(t *testing.T) {
	flex := NewFlex(Row).AddChild(0.5, "A")

	copia := flex.ChildrenCopy()
	copia[0].Content = "alterado"

	if flex.Children[0].Content != "A" {
		t.Error("ChildrenCopy devolveu o slice interno")
	}
}

func TestFlexSizeEChildCount(t *testing.T) {
	flex := NewFlex(Row).SetSize(33, 7).AddChild(0.5, "A").AddChild(0.5, "B")

	largura, altura := flex.Size()
	if largura != 33 || altura != 7 {
		t.Errorf("Size = (%d, %d), esperado (33, 7)", largura, altura)
	}
	if got := flex.ChildCount(); got != 2 {
		t.Errorf("ChildCount = %d, esperado 2", got)
	}
}

func TestFlexSemFilhos(t *testing.T) {
	flex := NewFlex(Row).SetSize(10, 2)
	if got := flex.Render(); got != "" {
		t.Errorf("Render = %q, esperado string vazia", got)
	}
}
