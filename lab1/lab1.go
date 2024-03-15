package lab1

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	FileTypeBin = "bin"
	FileTypeTxt = "txt"
)

func CalculateTextEntropy(input string) (float64, map[rune]float64) {
	frequency := make(map[rune]float64)
	for _, char := range input {
		frequency[char]++
	}

	entropy := 0.0
	for i := range frequency {
		probability := frequency[i] / float64(len(input))
		entropy -= probability * math.Log2(probability)
	}

	return entropy, frequency
}

func CalculateBinaryEntropy(input string) (float64, map[rune]float64) {
	frequency := make(map[rune]float64)
	for _, char := range input {
		frequency[char]++
	}

	p0 := frequency['0'] / float64(len(input))
	p1 := frequency['1'] / float64(len(input))

	entropy := 0.0
	if p0 > 0 {
		entropy -= p0 * math.Log2(p0)
	}
	if p1 > 0 {
		entropy -= p1 * math.Log2(p1)
	}

	return entropy, frequency
}

func ReadFile(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("no arguments provided")
	}

	filePath := args[0]
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileType := FileTypeTxt
	if len(args) > 1 {
		fileType = args[1]
	} else {
		fileType = strings.TrimPrefix(filepath.Ext(filePath), ".")
	}

	var data string
	switch fileType {
	case FileTypeBin:
		bytes, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		for _, b := range bytes {
			data += fmt.Sprintf("%08b", b)
		}
	case FileTypeTxt:
		bytes, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		data = string(bytes)
	default:
		return "", fmt.Errorf("invalid file type: %s. Only 'bin' or 'txt' are allowed", fileType)
	}

	return data, nil
}

func writeCSV(frequency map[rune]float64, filePath string) error {
	csvFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	err = csvWriter.Write([]string{"Character", "Frequency"})
	if err != nil {
		return err
	}

	for char, freq := range frequency {
		err = csvWriter.Write([]string{string(char), fmt.Sprintf("%f", freq)})
		if err != nil {
			return err
		}
	}

	return nil
}

func getInfoAmount(entropy float64, length int) float64 {
	infoAmount := entropy * float64(length)
	return infoAmount
}

func simulateError(data string, errorRate float64) (string, error) {
	var result []rune
	for _, bit := range data {
		randomNumber, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return "", err
		}
		if float64(randomNumber.Int64())/100 < errorRate {
			if bit == '0' {
				result = append(result, '1')
			} else {
				result = append(result, '0')
			}
		} else {
			result = append(result, bit)
		}
	}
	return string(result), nil
}

func EntropyStats(fileName string) {
	filePath := "./lab1/" + fileName
	fileType := strings.TrimPrefix(filepath.Ext(fileName), ".")
	data, err := ReadFile(filePath, fileType)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	var frequency map[rune]float64
	if fileType == FileTypeTxt {
		sentences := strings.Split(data, ",")
		latinSentence := strings.TrimSpace(sentences[1])
		cyrillicSentence := strings.TrimSpace(sentences[0])

		cEntropy, cFrequency := CalculateTextEntropy(cyrillicSentence)
		lEntropy, lFrequency := CalculateTextEntropy(latinSentence)

		for key, value := range lFrequency {
			cFrequency[key] = value
		}

		frequency = cFrequency

		fmt.Printf("Entropy Cyrillic, Latin Sentence: %f, %f\n", cEntropy, lEntropy)

		cyrillicInfoAmount := getInfoAmount(cEntropy, len(cyrillicSentence))
		latinInfoAmount := getInfoAmount(lEntropy, len(latinSentence))

		fmt.Printf("Info Amount Cyrillic, Latin Sentence: %f, %f\n", cyrillicInfoAmount, latinInfoAmount)
	} else {
		binEntropy, binFrequency := CalculateBinaryEntropy(data)
		frequency = binFrequency
		fmt.Printf("Entropy Bin: %f\n", binEntropy)

		for _, p := range []float64{0.1, 0.5, 1.0} {
			errorData, err := simulateError(data, p)
			if err != nil {
				log.Fatalf("Failed to simulate error: %s", err)
			}
			errorEntropy, _ := CalculateBinaryEntropy(errorData)
			fmt.Printf("Entropy for p=%v: %f\n", p, errorEntropy)
		}
	}

	err = writeCSV(frequency, "lab1/frequency.csv")
	if err != nil {
		log.Fatalf("Failed to write CSV: %s", err)
	}
}
