package api

import (
	"net/http"
	"strconv"
	"time"

	"gofer/pkg/db"
	"gofer/pkg/nextdate"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	idget := r.URL.Query().Get("id")
	if idget == "" {
		writeErr(w, http.StatusBadRequest, "не указан идентификатор задачи")
		return
	}

	id, err := strconv.ParseInt(idget, 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "неверный идентификатор")
		return
	}
	t, err := db.Get(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, " задача не найдена")
		return
	}
	// удаление задачи
	if t.Repeat == "" {
		err = db.Delete(id)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "не удалось удалить задачу")
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
		writeErr(w, http.StatusBadRequest, "не удалось определить дату")
		return
	}
	// обновление даты в бл
	err = db.UpdateDate(nextDate, t.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "не удалось обновить дату")
		return
	}
	writeJson(w, map[string]interface{}{})
}
