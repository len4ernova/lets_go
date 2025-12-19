package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"
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

	CREATE TABLE IF NOT EXISTS roles(
		role_id integer PRIMARY KEY AUTOINCREMENT,
		role_name varchar(25) NOT NULL DEFAULT "",
		role_description varchar(256) NOT NULL DEFAULT "",
		role_created_at varchar(8) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS hierarchy(
		id integer PRIMARY KEY AUTOINCREMENT,
		role_id integer NOT NULL,
		parent_id integer NOT NULL,
		FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES roles(role_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS permissions(
		permission_id integer PRIMARY KEY AUTOINCREMENT,
		permission_name varchar(10) NOT NULL DEFAULT "",
		permission_description varchar(256) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS resources(
		resource_id integer PRIMARY KEY AUTOINCREMENT,
		resource_name varchar(256) NOT NULL DEFAULT "",
		resource_parent integer NOT NULL DEFAULT 0,
		resource_description varchar(256) NOT NULL DEFAULT "",
		resource_create_at varchar(8) NOT NULL DEFAULT ""
	);
	CREATE TABLE IF NOT EXISTS access(
		id integer PRIMARY KEY AUTOINCREMENT,
		role_id integer NOT NULL,
		permission_id integer NOT NULL,
		resource_id integer NOT NULL,

		FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
		FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE,
		FOREIGN KEY (resource_id) REFERENCES resources(resource_id) ON DELETE CASCADE
	);
`

// const ind = `CREATE INDEX idx_snippets_created ON snippets(created);
// 			CREATE INDEX groups_glab_id ON works (glab_group_id);`

const patternTime = `20060102`

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if !checkExist(dsn) {
		//return db, create(db, schema+ind)
		return db, create(db, schema)
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
		fmt.Println("SCHEMa")
		return err
	}
	if err := insert(db); err != nil {
		fmt.Println("INSERT")
		return err
	}
	return nil
}

func insert(db *sql.DB) error {
	insertRole := `INSERT INTO roles (role_name, role_description, role_created_at)
		VALUES (:role_name, :role_description, :role_created_at);`

	type role struct {
		id          int
		name        string
		parentID    int
		description string
		createAt    string
	}
	roles := []role{
		{name: "admin", description: "", createAt: time.Now().Format(patternTime)},
		{name: "tech_manager", description: "", createAt: time.Now().Format(patternTime)},
		{name: "tech_lead", description: "", createAt: time.Now().Format(patternTime)},
		{name: "developer", description: "", createAt: time.Now().Format(patternTime)},
		{name: "platform-engineer", description: "", createAt: time.Now().Format(patternTime)},
		{name: "devops", description: "", createAt: time.Now().Format(patternTime)},
		{name: "junior-developer", description: "", createAt: time.Now().Format(patternTime)},
	}
	for _, r := range roles {
		_, err := db.Exec(insertRole,
			sql.Named("role_name", r.name),
			sql.Named("role_description", r.description),
			sql.Named("role_created_at", time.Now().Format(patternTime)))
		if err != nil {
			return err
		}
	}

	insertHierarchy := `INSERT INTO hierarchy (role_id, parent_id)
	VALUES (:role_id, :parent_id);`
	type hierarchy struct {
		role_id   int
		parent_id int
	}
	hierarchyRoles := []hierarchy{
		{role_id: 1, parent_id: 0},
		{role_id: 2, parent_id: 1},
		{role_id: 3, parent_id: 2},
		{role_id: 4, parent_id: 3},
		{role_id: 5, parent_id: 3},
		{role_id: 6, parent_id: 3},
		{role_id: 6, parent_id: 5},
		{role_id: 7, parent_id: 4},
		{role_id: 7, parent_id: 6},
	}
	for _, h := range hierarchyRoles {
		_, err := db.Exec(insertHierarchy,
			sql.Named("role_id", h.role_id),
			sql.Named("parent_id", h.parent_id))
		if err != nil {
			return err
		}
	}

	insertPermissions := `INSERT INTO permissions (permission_name, permission_description)
	VALUES (:permission_name, :permission_description);`

	permissions := map[string]string{"create": "create", "read": "read", "write": "write", "delete": "delete"}
	for right, dscr := range permissions {
		_, err := db.Exec(insertPermissions,
			sql.Named("permission_name", right),
			sql.Named("permission_description", dscr))
		if err != nil {
			return err
		}
	}
	inserResources := `INSERT INTO resources (resource_name, resource_parent, resource_description, resource_create_at)
	VALUES (:resource_name, :resource_parent, :resource_description, :resource_create_at);
	`

	type resource struct {
		name        string
		parentID    int
		description string
		createAt    string
	}
	resources := []resource{
		{name: "РБПО", parentID: 0, description: "процессы ГОСТ", createAt: time.Now().Format(patternTime)},
		{name: "Стенды", parentID: 0, description: "стенды", createAt: time.Now().Format(patternTime)},
		{name: "Хозяйство", parentID: 0, description: "прочее хозяйство", createAt: time.Now().Format(patternTime)},
		{name: "Планирование процессов разработки безопасного программного обеспечения", parentID: 1, description: "Планирование ...", createAt: time.Now().Format(patternTime)},
		{name: "Обучение сотрудников", parentID: 1, description: "Обучение ...", createAt: time.Now().Format(patternTime)},
		{name: "Формирование и предъявление требований безопасности к ПО", parentID: 1, description: "Формирование ...", createAt: time.Now().Format(patternTime)},
		{name: "Статический анализ", parentID: 1, description: "Статический анализ ...", createAt: time.Now().Format(patternTime)},
	}

	for _, resourse := range resources {
		_, err := db.Exec(inserResources,
			sql.Named("resource_name", resourse.name),
			sql.Named("resource_parent", resourse.parentID),
			sql.Named("resource_description", resourse.description),
			sql.Named("resource_create_at", resourse.createAt))
		if err != nil {
			return err
		}
	}

	insertAccess := `INSERT INTO access (role_id, permission_id, resource_id)
	VALUES (:role_id, :permission_id, :resource_id);
	`
	type access struct {
		role_id       int
		permission_id int
		resource_id   int
	}
	accessResources := []access{
		{role_id: 7, permission_id: 2, resource_id: 1}, // jun read
		{role_id: 7, permission_id: 3, resource_id: 1}, // jun +write
		{role_id: 4, permission_id: 1, resource_id: 1}, // developer  +create+ РБПО

		{role_id: 7, permission_id: 2, resource_id: 2}, // jun read
		{role_id: 7, permission_id: 3, resource_id: 2}, // jun +write
		{role_id: 5, permission_id: 1, resource_id: 2}, // developer  +create+ Стенды
		{role_id: 5, permission_id: 4, resource_id: 2}, // developer  +del+ Стенды

		{role_id: 3, permission_id: 4, resource_id: 1}, // tech_lead +del РБПО
	}
	for _, a := range accessResources {
		_, err := db.Exec(insertAccess,
			sql.Named("role_id", a.role_id),
			sql.Named("permission_id", a.permission_id),
			sql.Named("resource_id", a.resource_id))
		if err != nil {
			return err
		}
	}
	return nil
}
