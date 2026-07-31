package main

import (
	"log"
	"net/http"
)

func main() {

	FRONTEND := http.FileServer(http.Dir(http.Dir("./frontend"))
	http.Handle("/", fs)

	log.Println("Сервер фронтенда запущен на :3000)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
