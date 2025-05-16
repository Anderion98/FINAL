package server

import (
	"log"
	"net/http"
)

func Server(r http.Handler) error {
	port := ":7540"
	log.Println("Запуск сервера")
	log.Println("Сервер запущен на порту:", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("ошибка создания сервера", err)
		return err
	}
	return nil
}
