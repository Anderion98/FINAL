package server

import (
	"log"
	"net/http"
)

func Start() error {
	port := getPort()
	log.Printf("Запуск сервера на порту %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ошибка создания сервера", err)
		return err
	}
	return nil
}
