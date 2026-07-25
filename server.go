package main

import (
	"log"
	"net/http"
)

func main() {

	FRONTEND := http.FileServer(http.Dir("ДОБАВИТЬ СЮДА ПАПКУ ГДЕ ЛЕЖИТ FRONTEND"))
	http.Handle("/", fs)

	log.Println("Сервер фронтенда запущен на http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
