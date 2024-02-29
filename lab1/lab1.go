package lab1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func calculateEntropy(input string) (float64, map[rune]float64) {
	frequency := make(map[rune]float64)
	for _, char := range input {
		frequency[char]++
	}

	entropy := 0.0
	length := float64(len(input))
	for _, freq := range frequency {
		probability := freq / length
		entropy -= probability * math.Log2(probability)
	}

	return entropy, frequency
}

func readFile(filePath string, fileType string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var data string
	switch fileType {
	case "bin":
		info, _ := file.Stat()
		size := info.Size()
		bytes := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(bytes)
		for _, b := range bytes {
			data += fmt.Sprintf("%08b", b)
		}
	case "txt":
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data += scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid file type: %s. Only 'bin' or 'txt' are allowed", fileType)
	}

	return data, nil
}

func writeCSV(frequency map[rune]float64) error {
	csvfile, err := os.Create("./lab1/frequency.csv")
	if err != nil {
		return err
	}
	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	err = csvwriter.Write([]string{"Character", "Frequency"})
	if err != nil {
		return err
	}
	for char, freq := range frequency {
		err = csvwriter.Write([]string{string(char), fmt.Sprintf("%f", freq)})
		if err != nil {
			return err
		}
	}
	csvwriter.Flush()

	return nil
}

func Run(fullFileName string) {
	filePath := "./lab1/" + fullFileName
	fileType := strings.TrimPrefix(filepath.Ext(fullFileName), ".")
	data, err := readFile(filePath, fileType)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	entropy, frequency := calculateEntropy(data)
	fmt.Printf("Entropy: %f\n", entropy)

	err = writeCSV(frequency)
	if err != nil {
		log.Fatalf("Failed to write CSV: %s", err)
	}
}
