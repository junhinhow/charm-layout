package layout

import (
	"math"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// distribuir divide total unidades entre as fatias indicadas pelos pesos e
// garante que a soma das fatias seja exatamente igual a total.
//
// A divisão usa o método do maior resto: cada fatia recebe a parte inteira da
// sua fração e o que sobrar é distribuído entre as fatias com as maiores
// frações (empates resolvidos pelo índice menor). Assim nenhuma coluna é
// perdida e o container é preenchido por inteiro.
//
// Pesos nulos ou somando zero resultam numa divisão igualitária.
func distribuir(total int, pesos []float64) []int {
	n := len(pesos)
	fatias := make([]int, n)
	if n == 0 || total <= 0 {
		return fatias
	}

	soma := 0.0
	for _, p := range pesos {
		soma += p
	}

	if soma <= 0 {
		// Sem pesos utilizáveis: divide igualmente.
		base := total / n
		resto := total % n
		for i := range fatias {
			fatias[i] = base
			if i < resto {
				fatias[i]++
			}
		}
		return fatias
	}

	fracoes := make([]float64, n)
	usado := 0
	for i, p := range pesos {
		exato := float64(total) * p / soma
		inteiro := int(math.Floor(exato))
		fatias[i] = inteiro
		fracoes[i] = exato - float64(inteiro)
		usado += inteiro
	}

	ordem := make([]int, n)
	for i := range ordem {
		ordem[i] = i
	}
	sort.SliceStable(ordem, func(a, b int) bool {
		return fracoes[ordem[a]] > fracoes[ordem[b]]
	})

	resto := total - usado
	for i := 0; i < n && resto > 0; i++ {
		fatias[ordem[i]]++
		resto--
	}

	// Correção final para compensar eventuais arredondamentos de ponto
	// flutuante, mantendo a soma exatamente igual a total.
	for tentativa := 0; tentativa <= n; tentativa++ {
		atual := 0
		for _, f := range fatias {
			atual += f
		}
		if atual == total {
			break
		}
		if atual < total {
			fatias[ordem[0]]++
			continue
		}
		for _, idx := range ordem {
			if fatias[idx] > 0 {
				fatias[idx]--
				break
			}
		}
	}

	return fatias
}

// pesosIguais devolve n pesos iguais, úteis para dividir espaço de forma
// uniforme (colunas de uma grade, por exemplo).
func pesosIguais(n int) []float64 {
	pesos := make([]float64, n)
	for i := range pesos {
		pesos[i] = 1
	}
	return pesos
}

// limitarGap reduz o espaçamento para que a soma dos gaps caiba no espaço
// disponível, garantindo que o layout nunca ultrapasse o container.
//
// O valor devolvido é sempre não negativo e nunca maior que o gap solicitado.
func limitarGap(gap, espaco, n int) int {
	if n <= 1 || gap <= 0 {
		return 0
	}
	if maximo := espaco / (n - 1); gap > maximo {
		return maximo
	}
	return gap
}

// ajustarBloco normaliza um bloco já renderizado para que ocupe exatamente
// largura colunas e no máximo altura linhas.
//
// O lipgloss quebra o conteúdo dentro da largura informada, mas textos com
// caracteres largos (CJK, emoji) podem ultrapassar o limite em uma coluna. A
// normalização garante a geometria exata declarada pelo container.
func ajustarBloco(bloco string, largura, altura int) string {
	linhas := strings.Split(bloco, "\n")
	if altura > 0 && len(linhas) > altura {
		linhas = linhas[:altura]
	}

	for i := range linhas {
		if largura <= 0 {
			linhas[i] = ""
			continue
		}
		if lipgloss.Width(linhas[i]) > largura {
			linhas[i] = ansi.Truncate(linhas[i], largura, "")
		}
		if atual := lipgloss.Width(linhas[i]); atual < largura {
			linhas[i] += strings.Repeat(" ", largura-atual)
		}
	}

	return strings.Join(linhas, "\n")
}

// renderCelula renderiza o conteúdo dentro de uma célula com largura e altura
// exatas, truncando o que não couber.
func renderCelula(conteudo string, largura, altura int) string {
	if largura <= 0 {
		return ajustarBloco("", 0, altura)
	}

	estilo := lipgloss.NewStyle().Width(largura).Height(altura)
	return ajustarBloco(estilo.Render(conteudo), largura, altura)
}

// ratioValido informa se a razão pode ser usada no cálculo do layout.
func ratioValido(ratio float64) bool {
	if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
		return false
	}
	return ratio >= 0
}

// ratioSeguro devolve a razão utilizável: valores negativos, NaN ou infinitos
// são tratados como zero.
func ratioSeguro(ratio float64) float64 {
	if !ratioValido(ratio) {
		return 0
	}
	return ratio
}
