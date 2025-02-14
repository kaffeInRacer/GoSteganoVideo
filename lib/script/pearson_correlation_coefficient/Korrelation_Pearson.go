package pearson_correlation_coefficient

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

func CallCorrelationPearson(originalText, modifiedText, saveResultAt string) error {
	if len(originalText) < 2 || len(modifiedText) < 2 {
		return errors.New("originalText or modifiedText is too short")
	}

	modifiedParts := strings.Split(modifiedText, "\u0003")
	modifiedText = modifiedParts[0]

	var resultBuilder strings.Builder

	for length := 2; length <= len(originalText); length += 2 {
		subPlaintext := originalText[:length]
		subModified := modifiedText[:length]

		if len(subPlaintext) != len(subModified) {
			fmt.Printf("Skipping: Panjang subPlaintext (%d) dan subModified (%d) tidak sama\n", len(subPlaintext), len(subModified))
			continue
		}

		var sumX, sumY, sumXY, sumX2, sumY2 float64

		for i := 0; i < len(subPlaintext); i++ {
			x := float64(subPlaintext[i])
			y := float64(subModified[i])

			sumX += x
			sumY += y
			sumXY += x * y
			sumX2 += x * x
			sumY2 += y * y
		}

		n := float64(len(subPlaintext))
		numerator := (n * sumXY) - (sumX * sumY)
		denominator := math.Sqrt((n*sumX2 - math.Pow(sumX, 2)) * (n*sumY2 - math.Pow(sumY, 2)))

		if denominator == 0 {
			continue
		}

		correlation := numerator / denominator
		fmt.Printf("Karakter Original: %s , Karakter Modified: %s ,nilai korelasi: %.4f \n", subPlaintext, subModified, correlation)
		resultBuilder.WriteString(fmt.Sprintf("Panjang: %d, Koefisien Korelasi Pearson: %.4f\n", length, correlation))
	}

	path := fmt.Sprintf("%s/Koefiensi_Pearson.txt", saveResultAt)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(resultBuilder.String())
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
