package transposecipher

import (
	"strings"
)

// kriptografi transposisi kolom
// karna pembacaanya berdasarkan kolom yang disepakati

type transpose struct {
	plaintext  string
	ciphertext string
	padding    rune
}

const Padding = '\u00D7' //Unit Separator

func NewTranspose() *transpose {
	return &transpose{
		padding: Padding,
	}
}

// Encrypt is the transposition function that encodes the text by rearranging characters based on the key
func (t *transpose) Encrypt(text string, rowWidth, key int) string {
	textRune := []rune(text)
	colCount := (len(textRune) + rowWidth - 1) / rowWidth
	matrix := make([][]rune, colCount)

	// Initialize each row in the matrix
	for i := range matrix {
		matrix[i] = make([]rune, rowWidth)
	}

	for i := 0; i < colCount*rowWidth; i++ {
		col := i / rowWidth
		row := (key + i) % rowWidth
		if i < len(textRune) {
			matrix[col][row] = textRune[i]
			//fmt.Printf("karakter %c, baris: %d, kolom: %d, i: %d\n", textRune[i], row, col, i)
		} else {
			matrix[col][row] = t.padding
			//fmt.Printf("karakter %c, baris: %d, kolom: %d, i: %d\n", t.padding, row, col, i)
		}
	}

	var result strings.Builder
	for row := rowWidth - 1; row >= 0; row-- {
		for col := colCount - 1; col >= 0; col-- {
			result.WriteRune(matrix[col][row])
		}
	}
	return result.String()
}
