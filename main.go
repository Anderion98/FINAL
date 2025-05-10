package main

import (
	"gofer/package/api/nextdate"
	"gofer/package/db"
	"gofer/package/server"
	"net/http"

	"github.com/go-chi/chi"
)

func DirGetWeb() http.Handler {
	webDir := "./web"
	return http.FileServer(http.Dir(webDir))
}
func main() {
	db.New()
	r := chi.NewRouter()
	r.Handle("/*", DirGetWeb())
	r.Get("/api/nextdate", nextdate.NextDayHandler)
	server.Start(r)
}
