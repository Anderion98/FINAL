package api

import (
	"encoding/json"
	"fmt"
	"gofer/package/api/nextdate"
	"gofer/package/db"
	"log"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(nextdate.TimeFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	res, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, res)
}
func AddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t db.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if t.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.Add(&t)
	if err != nil {
		fmt.Printf("ошибка добавления задачи в базу: %v", err)
		writeJson(w, map[string]string{"error": "ошибка добавления задачи в базу"})
		return
	}
	result := map[string]any{
		"id": id,
	}

	writeJson(w, result)
}

/*func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}*/

// afterNow возвращает true, если date >= now (по дате, без времени)
func afterNow(now, date time.Time) bool {
	nowRounding := now.In(time.UTC).Truncate(24 * time.Hour)
	dateRounding := date.In(time.UTC).Truncate(24 * time.Hour)
	return dateRounding.After(nowRounding) || nowRounding.Equal(dateRounding)
}

func checkDate(t *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)

	if t.Date == "" {
		t.Date = time.Now().Format(nextdate.TimeFormat)
	}

	timeD, err := time.Parse(nextdate.TimeFormat, t.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}

	if !afterNow(now, timeD) {
		if t.Repeat == "" {
			// Если правило повторения нет - ставим сегодняшнюю дату
			t.Date = now.Format(nextdate.TimeFormat)
		} else {
			// Если правило есть - вычисляем следующую дату
			next, err := nextdate.NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
			}
			t.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("не удалось сериализовать ответ: %v", err)
		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
	}
}
