package objection

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Objection struct {
	AP               string `json:"ap"`
	ViolationNumber  string `json:"violation_number"`
	FirstNameDriver  string `json:"first_name_driver"`
	LastNameDriver   string `json:"last_name_driver"`
	MiddleNameDriver string `json:"middle_name_driver"`
	GenderDriver     string `json:"gender_driver"`
	ObjectionDate    string `json:"objection_date"`
	PublishDate      string `json:"publish_date"`
	DocumentType     int    `json:"document_type"`
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, ap, violationNumber, firstNameDriver, lastNameDriver, middleNameDriver, genderDriver,
	objectionDate, publishDate string, documentType int) {
	log.Println("Inserting objection record ...")
	insertObjectionSQL := `INSERT INTO objections(ap, violation_number, first_name_driver, last_name_driver, 
		middle_name_driver, gender_driver, objection_date, publish_date, document_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertObjectionSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, violationNumber, firstNameDriver, lastNameDriver, middleNameDriver, genderDriver,
		objectionDate, publishDate, documentType)
	if err != nil {
		log.Error(err)
	}
}

func Update(db *sql.DB, ap, publishDate string) {
	log.Println("Update violation record with record ...")
	updateObjection := `UPDATE objections set publish_date = ? WHERE ap = ?`
	statement, err := db.Prepare(updateObjection)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, ap, publishDate)
	if err != nil {
		log.Error(err)
	}
}

func basicGet(db *sql.DB, query string, arg ...interface{}) []*Objection {
	objections := []*Objection{}
	row, err := db.Query(query, arg...)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		objection := new(Objection)
		err = row.Scan(&objection.AP, &objection.ViolationNumber, &objection.FirstNameDriver, &objection.LastNameDriver,
			&objection.MiddleNameDriver, &objection.GenderDriver, &objection.ObjectionDate, &objection.PublishDate,
			&objection.DocumentType)
		if err != nil {
			log.Error(err)
		}
		objections = append(objections, objection)
	}

	return objections
}

func GetByAll(db *sql.DB) []*Objection {
	return basicGet(db, `SELECT ap, violation_number, first_name_driver, last_name_driver, middle_name_driver, 
	gender_driver, objection_date, publish_date, document_type FROM objections`)
}

func GetByViolationNumber(db *sql.DB, violationNumber string) *Objection {
	objections := basicGet(db, `SELECT ap, violation_number, first_name_driver, last_name_driver, middle_name_driver, 
	gender_driver, objection_date, publish_date, document_type FROM objections WHERE violation_number like ? `, violationNumber)
	if len(objections) != 0 {
		return objections[0]
	} else {
		return &Objection{}
	}
}

func GetByAP(db *sql.DB, ap string) *Objection {
	objections := basicGet(db, `SELECT ap, violation_number, first_name_driver, last_name_driver, middle_name_driver, 
	gender_driver, objection_date, publish_date, document_type FROM objections WHERE ap like ? `, ap)
	if len(objections) != 0 {
		return objections[0]
	} else {
		return &Objection{}
	}
}

func GetByObjectionDate(db *sql.DB, objetionDate string) []*Objection {
	return basicGet(db, `SELECT ap, violation_number, first_name_driver, last_name_driver, middle_name_driver, 
	gender_driver, objection_date, publish_date FROM objections WHERE objection_date like ?`, objetionDate)
}
