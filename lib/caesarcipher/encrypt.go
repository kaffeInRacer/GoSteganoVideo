package caesarcipher

import (
	"errors"
	"math"
	"slices"
)

type elements struct {
	U          []rune
	M          []rune
	V          []rune
	ciphertext []rune
	plaintext  []rune
	shifter    map[rune]rune
}

func NewCaesarCipher() *elements {
	return &elements{
		shifter: make(map[rune]rune),
	}
}

func indexOfRune(slice []rune, r rune) bool {
	for _, v := range slice {
		if v == r {
			return true
		}
	}
	return false
}

func validateInputs(K1 string, K2 int, message string) error {
	if len(K1) < 140 {
		return errors.New("K1 must have at least 140 characters")
	}
	if K1 == "" {
		return errors.New("K1 cannot be empty")
	}
	if K2 <= 0 {
		return errors.New("K2 must be greater than 0")
	}
	if message == "" {
		return errors.New("Message cannot be empty")
	}
	return nil
}

func (e *elements) Encrypt(Plaintext, K1 string, K2 int) (string, error) {
	if err := validateInputs(K1, K2, Plaintext); err != nil {
		return "", err
	}

	if len(e.U) == 0 {
		for _, char := range K1 {
			if !indexOfRune(e.U, char) {
				e.U = append(e.U, char)
				e.V = append(e.V, char)
			}
		}
	}

	e.shifter = make(map[rune]rune)
	lengthUnique := len(e.U)

	for i, char := range e.U {
		shiftedIndex := math.Mod(float64(i+K2), float64(lengthUnique))
		e.shifter[char] = e.U[int(shiftedIndex)]
	}

	for _, char := range Plaintext {
		if shiftedChar, exists := e.shifter[char]; exists {
			e.ciphertext = append(e.ciphertext, shiftedChar)
		} else {
			if !slices.Contains(e.M, char) {
				e.M = append(e.M, char)
			}
			for k := 32; k <= 126; k++ {
				newChar := rune(k)
				if !slices.Contains(e.V, newChar) && newChar != char {
					e.V = append(e.V, newChar)
					e.ciphertext = append(e.ciphertext, newChar)
					e.shifter[char] = newChar
					break
				}
			}
		}
	}

	return string(e.ciphertext) + "\u0003" + string(e.M), nil
}
