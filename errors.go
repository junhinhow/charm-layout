package layout

import "errors"

// Erros retornados pela validação de [Flex] e [Grid].
//
// Use errors.Is para testar um erro específico; as mensagens carregam o
// contexto do valor inválido que as originou.
var (
	// ErrInvalidDirection indica uma direção de flex desconhecida.
	ErrInvalidDirection = errors.New("layout: direção inválida")
	// ErrInvalidGap indica um espaçamento negativo.
	ErrInvalidGap = errors.New("layout: espaçamento inválido")
	// ErrInvalidSize indica largura ou altura negativa.
	ErrInvalidSize = errors.New("layout: tamanho inválido")
	// ErrInvalidRatio indica uma razão negativa, NaN ou infinita.
	ErrInvalidRatio = errors.New("layout: razão inválida")
	// ErrZeroRatioSum indica que todas as razões dos filhos somam zero.
	ErrZeroRatioSum = errors.New("layout: soma das razões é zero")
	// ErrInvalidColumns indica um número de colunas menor que 1.
	ErrInvalidColumns = errors.New("layout: número de colunas inválido")
)
