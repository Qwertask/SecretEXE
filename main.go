package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}

	defer file.Close()

	// Настраиваем логгер для записи в файл
	log.SetOutput(file)

	// Запуск HTTP сервера
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/main", mainHandler)

	http.HandleFunc("/eds", edsHandler)

	http.HandleFunc("/eds_check", eds_checkHandler)

	log.Println("Server started at :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
