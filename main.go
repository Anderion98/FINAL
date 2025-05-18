package main

import (
	"log"
	"net/http"

	"gofer/pkg/api"
	"gofer/pkg/db"
	"gofer/pkg/server"

	"github.com/go-chi/chi"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", api.NextDayHandler)
	r.Post("/api/task", api.AddTask)
	r.Get("/api/tasks", api.TasksHandler)
	r.Get("/api/task", api.GetTask)
	r.Put("/api/task", api.UpdateTask)
	r.Post("/api/task/done", api.DoneTask)
	r.Delete("/api/task", api.DeleteTask)

	server.Start(r)
}
