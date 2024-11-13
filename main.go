package main

import (
	"log"
	"net/http"
)

func main() {

	// Запуск HTTP сервера
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/main", uploadHandlerMain)

	http.HandleFunc("/eds", uploadHandlerEDS)

	http.HandleFunc("/eds_check", uploadHandlerEDS_CHECK)

	log.Println("Server started at :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
