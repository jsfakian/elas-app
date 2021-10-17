package officer

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

type Officer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rank      string `json:"rank"`
	Commander int    `json:"commander"`
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, firstName, lastName, rank string, commander int) {
	log.Println("Inserting officer record ...")
	insertOfficerSQL := `INSERT INTO officer(first_name, last_name, rank, commander) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertOfficerSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	res, err := statement.Exec(firstName, lastName, rank, commander)
	if err != nil {
		log.Error(err)
	}

	log.Info(res)
}

// We are passing db reference connection from main to our method with other parameters
func Update(db *sql.DB, firstName, lastName, rank string, commander int) {
	log.Println("Updating officer ...")
	updateOfficerSQL := `UPDATE officer set first_name = ?, last_name = ?, rank = ? WHERE commander = ?`
	statement, err := db.Prepare(updateOfficerSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(firstName, lastName, rank, commander)
	if err != nil {
		log.Error(err)
	}
}

func basicGet(db *sql.DB, query string, arg interface{}) *Officer {
	officer := new(Officer)
	row, err := db.Query(query, arg)
	if err != nil {
		log.Error(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&officer.FirstName, &officer.LastName, &officer.Rank)
		log.Info(officer)
		return officer
	}

	return nil
}

func GetByID(db *sql.DB, id int) *Officer {
	return basicGet(db, "SELECT * FROM officer WHERE id = ? ", id)
}

func GetByLastName(db *sql.DB, lastName string) *Officer {
	return basicGet(db, "SELECT * FROM officer WHERE last_name = ? ", lastName)
}

func GetByCommander(db *sql.DB, commander int) *Officer {
	return basicGet(db, "SELECT first_name, last_name, rank FROM officer WHERE commander = ?", commander)
}
