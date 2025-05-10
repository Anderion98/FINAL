package config

import (
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Port() string {
	port, exists := os.LookupEnv(EnvPort)
	if !exists || port == "" {
		log.Printf("Переменная окружения не задана, используется значение по умолчанию: %s", defaultPort)
		port = defaultPort
	}
	log.Printf("Запуск сервера на порту: %s", port)
	return ":" + port
}

func DbPath() string {
	storagePath, exists := os.LookupEnv(EnvDBFile)

	if !exists || storagePath == "" {
		storagePath = filepath.Join(EnvDbPath, DbName)
	} else {
		return EnvDBFile
	}

	return storagePath
}
