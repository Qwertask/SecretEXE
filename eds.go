package main

import (
	"log"
	"os/exec"
)

var certPath = "./certificates/certificate.pfx"
var filePath = "./secret/secret.exe"
var password = "12345678"

// removeSignatures удаляет все существующие подписи из файла
func removeSignatures(filePath string) error {
	// Удаление подписей с использованием signtool
	cmd := exec.Command("signtool", "remove", "/s", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to remove signatures: %v, output: %s", err, output)
		return err
	}

	log.Printf("Successfully removed signatures from file: %s\n", filePath)
	return nil
}

// signFile подписывает файл с использованием signtool
func signFile(filePath, certPath, password, timestampURL string) error {
	// Подписание файла с использованием метки времени
	cmd := exec.Command("signtool", "sign", "/f", certPath, "/p", password, "/fd", "SHA256", "/tr", timestampURL, "/td", "SHA256", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to sign file: %v, output: %s", err, output)
		return err
	}

	log.Printf("Successfully signed file: %s\n", filePath)
	return nil
}

// verifySignature проверяет подпись файла с использованием signtool
func verifySignature(filePath string) error {
	// Проверка подписи файла
	cmd := exec.Command("signtool", "verify", "/a", "/pa", filePath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("failed to verify signature: %v, output: %s", err, output)
		return err
	}
	log.Println(string(output))
	log.Printf("Successfully verified signature: %s\n", filePath)
	return nil
}
