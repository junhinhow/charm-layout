package layout

import "testing"

func TestGridColunasZeroNaoEntraEmPanic(t *testing.T) {
	// Regressão: Columns é público e pode ser zerado depois da criação,
	// o que causava "integer divide by zero" na renderização.
	grade := NewGrid(3)
	grade.Columns = 0
	grade.AddChild("a")

	if saida := grade.Render(); saida == "" {
		t.Error("Render deveria devolver uma saída não vazia")
	}
}

func TestGridColunasNegativasNaoEntraEmPanic(t *testing.T) {
	grade := NewGrid(3)
	grade.Columns = -5
	grade.AddChild("a")

	if saida := grade.Render(); saida == "" {
		t.Error("Render deveria devolver uma saída não vazia")
	}
}
