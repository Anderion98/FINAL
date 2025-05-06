package main

import (
	"gofer/package/server"
	"net/http"
)

func main() {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	server.Start()
}
