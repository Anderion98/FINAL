package db

import (
	"database/sql"
	"fmt"
	"gofer/package/other"
	"os"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	s := &Storage{}

	err := s.Init()
	if err != nil {
		fmt.Println(other.DefaultNameDB)
		return nil, err
	}

	return s, nil
}
func (s *Storage) Init() error {
	_, err := os.Stat(other.DefaultNameDB)
	var install bool
	if err != nil {
		install = true
	}
	if install {
		db, err := sql.Open("sqlite", other.DefaultNameDB)
		if err != nil {
			return err
		}
		_, err = db.Exec(other.Schema)
		if err != nil {
			return err
		}
		defer db.Close()
	}
	return nil
}
