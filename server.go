package main

import (
	"log"
	"net/http"
)

func main() {

	FRONTEND := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", FRONTEND)

	log.Println("Сервер фронтенда запущен на 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}

}
