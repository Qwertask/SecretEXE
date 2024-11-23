package main

import (
	"log"
	"os"
	"unicode/utf16"
)

const offset = int64(0x00001900)
const maxLen = 10000

func replaceTextInExe(newText string) error {
	// Открываем файл для чтения и записи в бинарном режиме
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("could not open file: %w", err)
		return err
	}
	defer file.Close()

	// Перемещаем указатель файла к заданному смещению
	_, err = file.Seek(offset, 0)
	if err != nil {
		log.Printf("could not seek to offset: %w", err)
		return err
	}

	// Конвертируем строку в массив UTF-16 кодовых единиц
	utf16Encoded := utf16.Encode([]rune(newText))

	// Создаем буфер нужного размера и заполняем его новыми данными в UTF-16
	buffer := make([]byte, 0, maxLen*2)
	for _, code := range utf16Encoded {
		buffer = append(buffer, byte(code), byte(code>>8))
	}

	// Заполняем буфер до максимальной длины пробелами с нулевыми байтами
	for len(buffer) < maxLen*2 {
		buffer = append(buffer, ' ', 0)
	}

	// Записываем буфер в файл
	_, err = file.Write(buffer[:maxLen*2])
	if err != nil {
		log.Printf("could not write to file: %w", err)
		return err
	}

	return nil
}
