package server

import (
	"gofer/package/config"
	"log"
	"net/http"
)

func Server(r http.Handler) error {
	port := config.Port()
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("ошибка создания сервера", err)
		return err
	}
	return nil
}
