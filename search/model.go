package search

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

type Violation struct {
	ID   int
	Code string
	Name string
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, code string, name string) {
	log.Println("Inserting violation record ...")
	insertViolationSQL := `INSERT INTO violations(code, name) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertViolationSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(code, name)
	if err != nil {
		log.Error(err)
	}
}

func GetByID(db *sql.DB, id int) Violation {
	viol := Violation{}
	row, err := db.Query("SELECT * FROM violations WHERE id = ? ", id)
	if err != nil {
		log.Error(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&viol.ID, &viol.Code, &viol.Name)
	}

	return viol
}

func GetByName(db *sql.DB, name string) Violation {
	viol := Violation{}
	row, err := db.Query("SELECT * FROM violations WHERE name = ? ", name)
	if err != nil {
		log.Error(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&viol.ID, &viol.Code, &viol.Name)
	}

	return viol
}
