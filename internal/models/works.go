package models

import (
	"database/sql"
	"errors"
	"time"
)

type Work struct {
	ID                 int
	GLGroupID          int
	GLGroupTitle       string
	GLGroupPath        string
	GLGroupCreatedAt   time.Time
	GLGroupDescription string
	Visible            bool
}
type WorkModel struct {
	DB *sql.DB
}

func (m *WorkModel) InsertWork(gl_id int, title string, path string, createdAt string, description string, visible bool) (int, error) {
	stmt := `INSERT INTO works (glab_group_id, glab_group_title, glab_group_path, glab_group_created_at, glab_group_description, visible)
VALUES(:gl_id, :title, :path, :created_at, :description, :visible)`

	result, err := m.DB.Exec(stmt,
		sql.Named("gl_id", gl_id),
		sql.Named("title", title),
		sql.Named("path", path),
		sql.Named("created_at", createdAt),
		sql.Named("description", description),
		sql.Named("visible", visible),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *WorkModel) GetWork(id int) (Work, error) {
	stmt := `SELECT work_id, glab_group_id, glab_group_title, glab_group_path, glab_group_created_at, glab_group_description, visible
	FROM works
	WHERE work_id = :id`
	row := m.DB.QueryRow(stmt, sql.Named("id", id))
	var w Work
	var visible int
	var t string
	err := row.Scan(&w.ID, &w.GLGroupID, &w.GLGroupTitle, &w.GLGroupPath, &t, &w.GLGroupDescription, &visible)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Work{}, ErrNoRecord
		} else {
			return Work{}, err
		}
	}
	w.GLGroupCreatedAt, err = time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return Work{}, err
	}
	if visible != 0 {
		w.Visible = true
	}
	// If everything went OK, then return the filled Snippet struct.
	return w, nil
}

// ListWorks - вернет весь список работ из БД (isAll=true), или 10 последних работ (isAll=false)
func (m *WorkModel) ListWorks(isAll bool) ([]Work, error) {
	var stmt string

	//fmt.Println("ListWorks 1")
	if isAll {
		stmt = `SELECT work_id, glab_group_id, glab_group_title, glab_group_path, glab_group_created_at, glab_group_description, visible
			FROM works
			ORDER BY glab_group_title ASC`
	} else {
		stmt = `SELECT work_id, glab_group_id, glab_group_title, glab_group_path, glab_group_created_at, glab_group_description, visible
			FROM works
			ORDER BY work_id DESC LIMIT 10`
	}
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	//fmt.Println("ListWorks 2")
	defer rows.Close()

	var works []Work

	for rows.Next() {
		var w Work
		var visible int
		var t string
		//	fmt.Println("ListWorks 3")
		err = rows.Scan(&w.ID, &w.GLGroupID, &w.GLGroupTitle, &w.GLGroupPath, &t, &w.GLGroupDescription, &visible)
		if err != nil {
			return nil, err
		}
		w.GLGroupCreatedAt, err = time.Parse(time.RFC3339Nano, t)
		if err != nil {
			return nil, err
		}
		if visible != 0 {
			w.Visible = true
		}
		//	fmt.Println("ListWorks 4", w)

		works = append(works, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the slice.
	return works, nil
}

// поиск по gitlab ID
func (m *WorkModel) GetWorkGLabID(id int) (Work, error) {
	stmt := `SELECT work_id, glab_group_id, glab_group_title, glab_group_path, glab_group_created_at, glab_group_description, visible
		FROM works
		WHERE glab_group_id = :glab_group_id`

	row := m.DB.QueryRow(stmt, sql.Named("glab_group_id", id))
	var w Work
	var visible int
	var t string
	err := row.Scan(&w.ID, &w.GLGroupID, &w.GLGroupTitle, &w.GLGroupPath, &t, &w.GLGroupDescription, &visible)
	if err != nil {
		return Work{}, err
	}
	w.GLGroupCreatedAt, err = time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return Work{}, err
	}
	if visible != 0 {
		w.Visible = true
	}
	return w, nil
}

func (m *WorkModel) UpdateWork(glab_group_id int, title string, path string, createdAt string, description string) error {
	stmt := `UPDATE works
	SET glab_group_title = :title, 
	glab_group_path = :path, 
	glab_group_created_at = :createdAt, 
	glab_group_description = :description
	WHERE glab_group_id = :glab_group_id`
	_, err := m.DB.Exec(stmt,
		sql.Named("title", title),
		sql.Named("path", path),
		sql.Named("createdAt", createdAt),
		sql.Named("description", description),
		sql.Named("glab_group_id", glab_group_id),
	)
	if err != nil {
		return err
	}
	return nil
}
