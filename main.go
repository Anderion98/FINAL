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
	r.Get("/api/task", api.GetTask)
	r.Put("/api/task", api.UpdateTask)
	r.Post("/api/task/done", api.DoneTask)
	r.Delete("/api/task", api.DeleteTask)

	server.Server(r)
}
