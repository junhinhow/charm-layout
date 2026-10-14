package layout

import (
	"errors"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestGridColunasZeroNaoEntraEmPanic(t *testing.T) {
	// Regressão: Columns é público e pode ser zerado depois da criação,
	// o que causava "integer divide by zero" na renderização.
	grade := NewGrid(3)
	grade.Columns = 0
	grade.AddChild("a")

	saida := grade.Render()
	if saida == "" {
		t.Error("Render deveria devolver uma saída não vazia")
	}

	if _, err := grade.RenderE(); !errors.Is(err, ErrInvalidColumns) {
		t.Fatalf("RenderE devolveu %v, esperado ErrInvalidColumns", err)
	}
}

func TestGridColunasNegativasNaoEntraEmPanic(t *testing.T) {
	grade := NewGrid(3)
	grade.Columns = -5
	grade.AddChild("a")

	if saida := grade.Render(); saida == "" {
		t.Error("Render deveria devolver uma saída não vazia")
	}
	if _, err := grade.RenderE(); !errors.Is(err, ErrInvalidColumns) {
		t.Fatalf("RenderE devolveu %v, esperado ErrInvalidColumns", err)
	}
}

func TestGridLarguraExata(t *testing.T) {
	grade := NewGrid(3).
		SetSize(80, 10).
		SetGap(1).
		AddChildren("A", "B", "C")

	saida := grade.Render()
	if got := lipgloss.Width(saida); got != 80 {
		t.Errorf("largura renderizada = %d, esperado 80", got)
	}
}

func TestGridColunasLarguraIgual(t *testing.T) {
	grade := NewGrid(3).
		SetSize(80, 6).
		SetGap(1).
		AddChildren("A", "B", "C")

	linha := primeiraLinha(grade.Render())

	if got := colunaDe(linha, "A"); got != 0 {
		t.Errorf("coluna da primeira célula = %d, esperado 0", got)
	}
	if got := colunaDe(linha, "B"); got != 27 {
		t.Errorf("coluna da segunda célula = %d, esperado 27", got)
	}
	if got := colunaDe(linha, "C"); got != 54 {
		t.Errorf("coluna da terceira célula = %d, esperado 54", got)
	}
}

func TestGridGapMaiorQueEspaco(t *testing.T) {
	// Regressão: com Gap=10, Columns=10 e width=80 a linha antiga renderizava
	// 91 colunas, estourando um terminal de 80 colunas.
	grade := NewGrid(10).
		SetSize(80, 10).
		SetGap(10).
		AddChildren("A", "B", "C", "D", "E", "F", "G", "H", "I", "J")

	saida := grade.Render()
	if got := lipgloss.Width(saida); got > 80 {
		t.Errorf("largura renderizada = %d, esperado no máximo 80", got)
	}
}

func TestGridDistribuiTodasAsColunas(t *testing.T) {
	// Uma divisão inteira simples perderia colunas: 10 colunas em 80 colunas
	// com gap 1 dão 7 colunas de largura 7 e sobra 1.
	grade := NewGrid(10).
		SetSize(80, 10).
		SetGap(1).
		AddChildren("A", "B", "C", "D", "E", "F", "G", "H", "I", "J")

	saida := grade.Render()
	if got := lipgloss.Width(saida); got != 80 {
		t.Errorf("largura renderizada = %d, esperado 80", got)
	}
}

func TestGridNumeroDeLinhas(t *testing.T) {
	grade := NewGrid(2).
		SetSize(20, 3).
		AddChildren("A", "B", "C", "D", "E")

	linhas := strings.Split(grade.Render(), "\n")
	if len(linhas) != 3 {
		t.Errorf("número de linhas = %d, esperado 3", len(linhas))
	}
}

func TestGridAlturaExata(t *testing.T) {
	grade := NewGrid(2).
		SetSize(20, 6).
		AddChildren("A", "B", "C", "D")

	saida := grade.Render()
	if got := lipgloss.Height(saida); got != 6 {
		t.Errorf("altura renderizada = %d, esperado 6", got)
	}
}

func TestGridConteudoLongoNaoEstoura(t *testing.T) {
	grade := NewGrid(2).
		SetSize(10, 4).
		AddChildren("conteúdo bem maior que a célula", "B")

	saida := grade.Render()
	if got := lipgloss.Width(saida); got != 10 {
		t.Errorf("largura renderizada = %d, esperado 10", got)
	}
}

func TestGridSetSizeNegativo(t *testing.T) {
	grade := NewGrid(2).SetSize(-1, 5).AddChild("A")
	if _, err := grade.RenderE(); !errors.Is(err, ErrInvalidSize) {
		t.Fatalf("RenderE devolveu %v, esperado ErrInvalidSize", err)
	}
}

func TestGridSetGapNegativo(t *testing.T) {
	grade := NewGrid(2).SetGap(-2).AddChild("A")
	if _, err := grade.RenderE(); !errors.Is(err, ErrInvalidGap) {
		t.Fatalf("RenderE devolveu %v, esperado ErrInvalidGap", err)
	}
}

func TestGridNewGridColunasInvalidas(t *testing.T) {
	if _, err := NewGridE(0); !errors.Is(err, ErrInvalidColumns) {
		t.Fatalf("NewGridE(0) devolveu %v, esperado ErrInvalidColumns", err)
	}
	if _, err := NewGridE(3); err != nil {
		t.Fatalf("NewGridE(3) devolveu %v, esperado nil", err)
	}
}

func TestGridChildrenCopy(t *testing.T) {
	grade := NewGrid(2).AddChildren("A", "B")

	copia := grade.ChildrenCopy()
	copia[0] = "alterado"

	if grade.Children[0] != "A" {
		t.Error("ChildrenCopy devolveu o slice interno")
	}
}

func TestGridSemFilhos(t *testing.T) {
	grade := NewGrid(2).SetSize(10, 2)
	if got := grade.Render(); got != "" {
		t.Errorf("Render = %q, esperado string vazia", got)
	}
}
