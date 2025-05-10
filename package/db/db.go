package db

import (
	"database/sql"
	"gofer/package/config"
	"os"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func DataBase() (*Storage, error) {
	s := &Storage{}

	err := s.Init()
	if err != nil {
		return nil, err
	}

	return s, nil
}
func (s *Storage) Init() error {
	_, err := os.Stat(config.DbName)
	var install bool
	if err != nil {
		install = true
	}
	if install {
		db, err := sql.Open("sqlite", config.DbName)
		if err != nil {
			return err
		}
		_, err = db.Exec(config.Schema)
		if err != nil {
			return err
		}
		defer db.Close()
	}
	return nil
}
