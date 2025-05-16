package api

import (
	"fmt"
	"gofer/package/api/nextdate"
	"gofer/package/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	idget := r.URL.Query().Get("id")
	if idget == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор задачи"})
		return
	}

	id, err := strconv.ParseInt(idget, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "некорректный идентификатор задачи"})
		return
	}
	t, err := db.Get(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"})
		return
	}
	// удаление задачи
	if t.Repeat == "" {
		err = db.Delete(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]interface{}{})
		return
	}
	// переодичская задача, обнуляем дату
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	nextDate, err := nextdate.NextDate(now, t.Date, t.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("дата не получена: %v", err)})
		return
	}
	// обновление даты в бл
	err = db.UpdateDate(nextDate, t.ID)
	if err != nil {
		log.Println(66)
		writeJson(w, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, map[string]interface{}{})
}
