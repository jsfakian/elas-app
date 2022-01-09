package violation

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

type Violation struct {
	AP                 string `json:"ap"`
	AT                 string `json:"at"`
	ViolationNumber    string `json:"violation_number"`
	RegistrationNumber string `json:"registration_number"`
	FirstNameOwner     string `json:"first_name_owner"`
	MiddleNameOwner    string `json:"middle_name_owner"`
	LastNameOwner      string `json:"last_name_owner"`
	AddressOwner       string `json:"address_owner"`
	DocumentType       int    `json:"document_type"`
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, ap, at, violationNumber, registrationNumber, firstNameOwner,
	lastNameOwner, middleNameOwner, addressOwner string, documentType int) {
	log.Println("Inserting violation record ...")
	insertViolationSQL := `INSERT INTO violations(ap, at, violation_number,
		registration_number, first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertViolationSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, at, violationNumber, registrationNumber,
		firstNameOwner, middleNameOwner, lastNameOwner, addressOwner, documentType)
	if err != nil {
		log.Error(err)
	}
}

func UpdatePublishDate(db *sql.DB, ap, publishDate string) {
	log.Println("Update violation record with record ...")
	updateViolation := `UPDATE violations set publish_date = ? WHERE ap = ?`
	statement, err := db.Prepare(updateViolation)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, publishDate)
	if err != nil {
		log.Error(err)
	}
}

func basicGet(db *sql.DB, query string, arg ...interface{}) []*Violation {
	violations := []*Violation{}

	row, err := db.Query(query, arg...)
	if err != nil {
		log.Error("Failed query: ", err)
		return nil
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		violation := new(Violation)
		err = row.Scan(&violation.AP, &violation.AT, &violation.ViolationNumber, &violation.RegistrationNumber,
			&violation.FirstNameOwner, &violation.MiddleNameOwner, &violation.LastNameOwner, &violation.AddressOwner,
			&violation.DocumentType)
		if err != nil {
			log.Error("Failed scanning: ", err)
		}
		log.Info(violation)
		violations = append(violations, violation)
	}

	return violations
}

func GetByAll(db *sql.DB) []*Violation {
	return basicGet(db, `SELECT ap, at, violation_number, registration_number, 
	first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type FROM violations`)
}

func GetByAP(db *sql.DB, ap string) *Violation {
	log.Info("ap: ", ap)
	return basicGet(db, `SELECT ap, at, violation_number, registration_number, 
	first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type FROM violations 
	WHERE ap like ? `, ap)[0]
}

func GetByViolationNumber(db *sql.DB, violationNumber string) *Violation {
	return basicGet(db, `SELECT ap, at, violation_number, registration_number, 
	first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type FROM violations 
	WHERE violation_number like ? `, violationNumber)[0]
}

func GetByDriver(db *sql.DB, firstName, lastName string) []*Violation {
	return basicGet(db, `SELECT ap, at, violation_number, registration_number, 
	first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type FROM violations 
	WHERE first_name_driver like ? and last_name_driver like ?`, firstName, lastName)
}

func GetByOwner(db *sql.DB, firstName, lastName string) []*Violation {
	return basicGet(db, `SELECT ap, at, violation_number, registration_number, 
	first_name_owner, middle_name_owner, last_name_owner, address_owner, document_type FROM violations 
	WHERE first_name_owner like ? and last_name_owner like ?`, firstName, lastName)
}
