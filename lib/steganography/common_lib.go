package steganography

const emmit = '\u0005' // Padding Character

func stringToBinary(str string) []uint8 {
	var outputRes []uint8
	bytes := []byte(str)
	for _, char := range bytes {
		for i := 7; i >= 0; i-- {
			bit := (char >> uint(i)) & 1
			outputRes = append(outputRes, bit)
		}
	}
	return outputRes
}
func binaryToString(bits []uint8) (string, bool) {
	var outputRes []byte
	for i := 0; i < len(bits); i += 8 {
		var byteVal uint8
		for j := 0; j < 8; j++ {
			if i+j < len(bits) {
				byteVal = (byteVal << 1) | bits[i+j]
			}
		}
		if byteVal == emmit {
			return string(outputRes), true
		}
		outputRes = append(outputRes, byteVal)
	}
	return string(outputRes), false
}
