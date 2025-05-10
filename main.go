package main

import (
	"gofer/package/db"
	"gofer/package/server"
	"net/http"
)

func main() {
	db.New()
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	server.Start()
}
