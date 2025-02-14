package caesarcipher

import (
	"slices"
	"strings"
)

func (e *elements) Decrypt(Ciphertext, K1 string, K2 int) (string, error) {
	if err := validateInputs(K1, K2, Ciphertext); err != nil {
		return "", err
	}
	newElement := new(elements)
	newElement.shifter = make(map[rune]rune)

	parts := strings.Split(Ciphertext, "\u0003")
	Ciphertext = parts[0]
	newElement.M = []rune(parts[1])
	if len(parts) > 1 {
		newElement.M = []rune(parts[1])
	}

	if len(newElement.U) == 0 {
		for _, char := range K1 {
			if !indexOfRune(newElement.U, char) {
				newElement.U = append(newElement.U, char)
				newElement.V = append(newElement.V, char)
			}
		}
	}

	lengthUnique := len(newElement.U)

	for i, char := range newElement.U {
		shiftedIndex := (i + lengthUnique - (K2 % lengthUnique)) % lengthUnique
		newElement.shifter[char] = newElement.U[shiftedIndex]
	}

	for _, ms := range newElement.M {
		for k := 32; k <= 126; k++ {
			newChar := rune(k)
			if !slices.Contains(newElement.V, newChar) && newChar != ms {
				newElement.V = append(newElement.V, newChar)
				newElement.shifter[newChar] = ms
				break
			}
		}
	}

	for _, char := range Ciphertext {
		if originalChar, exists := newElement.shifter[char]; exists {
			newElement.plaintext = append(newElement.plaintext, originalChar)
		}
	}

	return string(newElement.plaintext), nil
}
