package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gofer/pkg/db"
	"gofer/pkg/nextdate"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func writeErr(w http.ResponseWriter, errcode int, text string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errcode)
	json.NewEncoder(w).Encode(map[string]any{"error": text})
}

// не использовал writeJson ввиду несовпадения типов
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
		writeErr(w, http.StatusBadRequest, "JSON ошибка")
		return
	}

	if t.Title == "" {
		writeErr(w, http.StatusBadRequest, "ошибка в поле Title")
		return
	}

	if err := checkDate(&t); err != nil {
		writeErr(w, http.StatusBadRequest, "не удалось проверить дату")
		return
	}

	id, err := db.Add(&t)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "не удалось добавить задачу")
		return
	}
	result := map[string]any{
		"id": id,
	}

	writeJson(w, result)
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "не удалось получить задачи")
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	// поле ID не совпадает с типом поля в JSON, в ручную меняем тип и преобразуем в мапу
	respTasks := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		taskMap := map[string]string{
			"id":      fmt.Sprintf("%d", t.ID),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		}
		respTasks = append(respTasks, taskMap)
	}

	writeJson(w, map[string]interface{}{"tasks": respTasks})
}
func GetTask(w http.ResponseWriter, r *http.Request) {
	idget := r.URL.Query().Get("id")
	if idget == "" {
		writeErr(w, http.StatusBadRequest, "id пустое")
		return
	}
	id, err := strconv.ParseInt(idget, 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id")
		return
	}

	task, err := db.Get(id)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "задача не найдена")
		return
	}

	//ответ в виде JSON не совпадает с типом поля, меняем тип поля вручную
	writeJson(w, map[string]string{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErr(w, http.StatusBadRequest, "JSON ошибка")
		return
	}

	if input.ID == "" {
		writeErr(w, http.StatusBadRequest, "id пустое")
		return
	}
	// парсинг ID вручную
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id ")
		return
	}

	if input.Title == "" {
		writeErr(w, http.StatusBadRequest, " ошибка в поле Title")
		return
	}

	// заполняем структуру
	task := &db.Task{
		ID:      id,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	// проверка даты
	if err := checkDate(task); err != nil {
		writeErr(w, http.StatusBadRequest, "не удалось проверить дату")
		return
	}

	// обновление задачи в БД
	if err := db.Update(task); err != nil {
		writeErr(w, http.StatusNotFound, "не удалось обновить задачу")
		return
	}

	writeJson(w, map[string]interface{}{})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idget := r.URL.Query().Get("id")
	if idget == "" {
		writeErr(w, http.StatusBadRequest, "id пустое")
		return
	}
	id, err := strconv.ParseInt(idget, 10, 64)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id")
		return
	}

	err = db.Delete(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "не удалось удалить задачу")
		return
	}

	writeJson(w, map[string]interface{}{})
}

// afterNow не стал импортировать из nextdate, создал здесь свое
func afterNow(now, date time.Time) bool {
	nowRounding := now.Truncate(24 * time.Hour)
	dateRounding := date.Truncate(24 * time.Hour)
	return dateRounding.After(nowRounding) || nowRounding.Equal(dateRounding)
}

func checkDate(t *db.Task) error {
	r := time.Now()
	now := time.Date(r.Year(), r.Month(), r.Day(), 0, 0, 0, 0, r.Location())
	if t.Date == "" {
		t.Date = time.Now().Format(nextdate.TimeFormat)
	}

	timeD, err := time.Parse(nextdate.TimeFormat, t.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}

	if !afterNow(now, timeD) {

		if t.Repeat == "" {

			// если правила нет ставим текущую дату в нужном формате
			t.Date = now.Format(nextdate.TimeFormat)

		} else {

			// вычисляем следующую дату
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
	}
}
