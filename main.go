package main

import (
	"fmt"
	"os"
	"unicode"

	"main.go/redact"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Пожалуйста введите два аргумента")
		return
	}
	if len(os.Args[1]) <= 4 || len(os.Args[2]) <= 4 || os.Args[1][len(os.Args[1])-4:] != ".txt" || os.Args[2][len(os.Args[2])-4:] != ".txt" {
		fmt.Println("Пожалуйста, введите корректные формат и название файла (.txt)")
		return
	}

	data := os.Args[1]
	input, err := os.ReadFile(data)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла %s: %v\n", data, err)
		return
	}
	res := string(input)
	for _, ch := range res {
		if ch < unicode.MaxASCII {
		} else {
			fmt.Println("Введите в sample.txt только ASCII символы")
			return
		}
	}
	res = redact.RedactFile(res)
	result := os.Args[2]
	output, err := os.Create(result)
	if err != nil {
		fmt.Printf("Ошибка при создании файла %s: %v\n", result, err)
		return
	}
	defer output.Close()
	_, err = output.WriteString(res)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл %s: %v\n", result, err)
		return
	}

	fmt.Printf("Результат успешно записан в %s\n", result)
}
