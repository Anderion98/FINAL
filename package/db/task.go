package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Add(task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT * 
		 FROM scheduler
		 ORDER BY date ASC
		 LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := new(Task)
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}
		tasks = append(tasks, t)
	}
	/*if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %w", err)
	}*/

	return tasks, nil
}
func InitDB(dateb *sql.DB) {
	db = dateb
}
