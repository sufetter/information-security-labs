package lab2

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"infsec/lab1"
	"log"
	"math"
	"strings"
)

func encodeToBase64(input string) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	data := []byte(input)

	for i := 0; i < len(data); i += 3 {
		val := uint32(data[i]) << 16
		if i+1 < len(data) {
			val |= uint32(data[i+1]) << 8
		}
		if i+2 < len(data) {
			val |= uint32(data[i+2])
		}

		for j := 18; j >= 0; j -= 6 {
			if j == 6 && len(data)-i < 3 {
				result.WriteByte('=')
			} else if j == 0 && len(data)-i < 2 {
				result.WriteByte('=')
			} else {
				result.WriteByte(base64Chars[(val>>uint(j))&0x3F])
			}
		}
	}

	missing := len(data) % 3
	if missing > 0 {
		for i := 0; i < 3-missing; i++ {
			result.WriteByte('=')
		}
	}

	return result.String()

	//То же самое, но с использованием
	//стандартной библиотеки

	//return base64.StdEncoding.EncodeToString([]byte(data))
}

func informationRedundancy(input string, alphabetSize int, entropy ...float64) (float64, error) {
	var entropyVal float64
	if len(entropy) > 0 {
		entropyVal = entropy[0]
	} else {
		entropyVal, _ = lab1.CalculateTextEntropy(input)
	}

	maxEntropy := math.Log2(float64(alphabetSize))

	return maxEntropy - entropyVal, nil
}

func xorBuffers(a, b []byte) []byte {
	//Насколько мне известно, в Go нет стандартных функций
	//для подобных операций, поэтому результат моей функции
	//сравнить не с чем, но я уверен в его корректности

	if len(a) < len(b) {
		a = append(a, b[len(a):]...)
	} else if len(b) < len(a) {
		b = append(b, a[len(b):]...)
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}

	return result
}

func InformationStats(fileName string) {
	data, err := lab1.ReadFile("./lab2/" + fileName)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	dataBase64 := encodeToBase64(data)

	//fmt.Printf("Base64 version: %v\n", dataBase64)

	latinEntropy, _ := lab1.CalculateTextEntropy(data)
	base64Entropy, _ := lab1.CalculateTextEntropy(dataBase64)
	fmt.Printf("Latin entropy: %f\nBase64 entropy: %f\n", latinEntropy, base64Entropy)

	latinRedundancy, _ := informationRedundancy(data, 26, latinEntropy)
	base64Redundancy, _ := informationRedundancy(dataBase64, 64, base64Entropy)
	fmt.Printf("Latin redundancy: %f\nBase64 redundancy: %f\n", latinRedundancy, base64Redundancy)

	base64Name := encodeToBase64("Daniil")
	fmt.Printf("Is the result of my encoder equal to the result from std encoder? \n%v\n", bytes.Equal([]byte(base64Name), []byte(base64.StdEncoding.EncodeToString([]byte("Daniil")))))
	base64Surname := encodeToBase64("Lobunets")
	fmt.Printf("XOR result ASCII: %x\n", xorBuffers([]byte("Lobunets"), []byte("Daniil")))
	fmt.Printf("XOR result Base64: %x\n", xorBuffers([]byte(base64Surname), []byte(encodeToBase64(base64Name))))
}
