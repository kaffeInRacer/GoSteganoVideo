package transposecipher

import (
	"strings"
)

func (t *transpose) Decrypt(text string, rowWidth, key int) string {
	textRune := []rune(text)
	colCount := (len(textRune) + rowWidth - 1) / rowWidth
	matrix := make([][]rune, colCount)
	for i := range matrix {
		matrix[i] = make([]rune, rowWidth)
	}

	index := 0
	for row := rowWidth - 1; row >= 0; row-- {
		for col := colCount - 1; col >= 0; col-- {
			matrix[col][row] = textRune[index]
			index++
		}
	}

	var result strings.Builder
	for i := range colCount * rowWidth {
		col := i / rowWidth
		row := (key + i) % rowWidth
		if matrix[col][row] != t.padding {
			result.WriteRune(matrix[col][row])
		}
	}

	return result.String()
}
