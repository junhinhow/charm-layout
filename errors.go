package layout

import "errors"

// Erros retornados pela validação de [Flex] e [Grid].
//
// Use errors.Is para testar um erro específico; as mensagens carregam o
// contexto do valor inválido que as originou.
var (
	// ErrDirecaoInvalida indica uma direção de flex desconhecida.
	ErrDirecaoInvalida = errors.New("layout: direção inválida")
	// ErrGapInvalido indica um espaçamento negativo.
	ErrGapInvalido = errors.New("layout: espaçamento inválido")
	// ErrTamanhoInvalido indica largura ou altura negativa.
	ErrTamanhoInvalido = errors.New("layout: tamanho inválido")
	// ErrRatioInvalido indica uma razão negativa, NaN ou infinita.
	ErrRatioInvalido = errors.New("layout: razão inválida")
	// ErrRazaoSomaZero indica que todas as razões dos filhos somam zero.
	ErrRazaoSomaZero = errors.New("layout: soma das razões é zero")
	// ErrColunasInvalidas indica um número de colunas menor que 1.
	ErrColunasInvalidas = errors.New("layout: número de colunas inválido")
)
