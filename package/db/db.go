package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbfile string) error {
	var err error
	db, err = sql.Open("sqlite", dbfile)
	if err != nil {
		return err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	db.Close()
}
