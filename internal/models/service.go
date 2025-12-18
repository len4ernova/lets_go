package models

import (
	"database/sql"
	"os"
)

const schema = `
	CREATE TABLE IF NOT EXISTS works(
		work_id integer PRIMARY KEY AUTOINCREMENT,
		glab_group_id integer NOT NULL DEFAULT 0,
		glab_group_title varchar(20) NOT NULL DEFAULT "",
		glab_group_path varchar(256) NOT NULL DEFAULT "",
		glab_group_created_at varchar(8) NOT NULL DEFAULT "",
		glab_group_description varchar(256) NOT NULL DEFAULT "",
		visible integer NOT NULL DEFAULT 1
	);
	`

const ind = `CREATE INDEX idx_snippets_created ON snippets(created);
			CREATE INDEX groups_glab_id ON works (glab_group_id);`

const patternTime = `20060102`

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if !checkExist(dsn) {
		return db, create(db, schema+ind)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// checkExist - проверка существования файла БД.
func checkExist(dbFile string) bool {
	_, err := os.Stat(dbFile)
	if err != nil {
		return false
	}

	return true
}

// create - создание таблиц.
func create(db *sql.DB, schema string) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
