package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const maxFileSize = 200 * 1024 // 400 кб

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("input")
	log.Println("Received input:", input)

	if len(input) > 10000 {
		http.Error(w, "Input string too long", http.StatusBadRequest)
		return
	}

	// Внесение изменений в exe файл
	err := replaceTextInExe(input)
	if err != nil {
		http.Error(w, "Error replacing text in exe file", http.StatusBadRequest)
		return
	}

	// Подписание измененного файла
	timestampURL := "http://timestamp.digicert.com"
	err = signFile(filePath, certPath, password, timestampURL)
	if err != nil {
		log.Printf("Error signing file: %v\n", err)
		//http.Error(w, "Error signing file", http.StatusInternalServerError)
	}

	// Открытие файла для отправки клиенту
	appFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Could not open file", http.StatusInternalServerError)
		return
	}
	defer appFile.Close()

	// Отправка файла клиенту
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, appFile)
	err = removeSignatures(filePath)
	if err != nil {
		log.Printf("Error removing signatures: %v\n", err)
		//http.Error(w, "Error removing signatures", http.StatusInternalServerError)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlFile, err := os.Open("./src/main.html")
	if err != nil {
		http.Error(w, "Could not open HTML file", http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()
	htmlData, err := ioutil.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, "Could not read HTML file", http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}

func edsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlFile, err := os.Open("./src/eds.html")
	if err != nil {
		http.Error(w, "Could not open HTML file", http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()
	htmlData, err := ioutil.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, "Could not read HTML file", http.StatusInternalServerError)
		return
	}
	w.Write(htmlData)
}

func eds_checkHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Ошибка получения файла:", err)
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", handler.Filename)
	if err != nil {
		log.Println("Ошибка создания временного файла:", err)
		http.Error(w, "Ошибка создания временного файла", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Println("Ошибка копирования файла:", err)
		http.Error(w, "Ошибка копирования файла", http.StatusInternalServerError)
		return
	}

	tempFile.Close()

	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	err = verifySignature(tempFile.Name())
	if err != nil {
		log.Println("Error verifying signature:", err)
		response.Status = "invalid"
		response.Message = "Подпись недействительна"
	} else {
		log.Printf("Загружен файл: %s\n", handler.Filename)
		response.Status = "valid"
		response.Message = "Подпись действительна"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
