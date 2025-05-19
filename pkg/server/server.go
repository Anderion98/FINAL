package server

import (
	"log"
	"net/http"
	"os"
)

const DefaultPort = "7540"

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = DefaultPort
		log.Printf("Переменная окружения не задана, порт по умолчанию: %v", DefaultPort)
	}
	return ":" + port
}
func Start(r http.Handler) error {
	port := getPort()
	log.Printf("Запуск сервера на порту %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("ошибка запуска сервера :", err)
		return err
	}
	return nil
}
