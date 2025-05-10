package server

import (
	"log"
	"net/http"
)

func Start(r http.Handler) error {
	port := getPort()
	log.Printf("Запуск сервера на порту %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("ошибка создания сервера", err)
		return err
	}
	return nil
}
