package officer

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

type Officer struct {
	ID        int
	FirstName string
	LastName  string
	Rank      string
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, firstName, lastName, rank string) {
	log.Println("Inserting violation record ...")
	insertViolationSQL := `INSERT INTO officer(first_name, last_name, rank) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertViolationSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(firstName, lastName, rank)
	if err != nil {
		log.Error(err)
	}
}

func GetByID(db *sql.DB, id int) Officer {
	officer := Officer{}
	row, err := db.Query("SELECT * FROM officer WHERE id = ? ", id)
	if err != nil {
		log.Error(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&officer.ID, &officer.FirstName, &officer.LastName, &officer.Rank)
	}

	return officer
}

func GetByName(db *sql.DB, lastName string) Officer {
	officer := Officer{}
	row, err := db.Query("SELECT * FROM officer WHERE last_name = ? ", lastName)
	if err != nil {
		log.Error(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&officer.ID, &officer.FirstName, &officer.LastName, &officer.Rank)
	}

	return officer
}
