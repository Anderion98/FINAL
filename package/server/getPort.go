package server

import (
	"gofer/package/other"
	"log"
	"os"
)

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = other.DefaultPort
		log.Printf("Переменная окружения не задана, используется значение по умолчанию: %v", other.DefaultPort)
		return ":" + port
	}
	log.Printf("Значение переменной окружения: %v", port)
	return ":" + port
}
