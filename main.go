package main

import (
	"final/package/api"
	"gofer/package/db"
	"gofer/package/server"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func DirGetWeb() http.Handler {
	webDir := "./web"
	return http.FileServer(http.Dir(webDir))
}
func main() {
	dbfile := "scheduler.db"
	if err := db.Init(dbfile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Handle("/*", DirGetWeb())
	r.Get("/api/nextdate", api.NextDayHandler)
	r.Post("/api/task", api.AddTask)
	r.Get("/api/tasks", api.TasksHandler)
	server.Server(r)
}
