package decision

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Decision struct {
	AP              string `json:"ap"`
	ViolationNumber string `json:"violation_number"`
	PublishDate     string `json:"publish_date"`
}

// We are passing db reference connection from main to our method with other parameters
func Insert(db *sql.DB, ap, violationNumber, publishDate string) {
	log.Println("Inserting decision record ...")
	insertDecisionSQL := `INSERT INTO decisions(ap, violation_number, publish_date) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertDecisionSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, violationNumber, publishDate)
	if err != nil {
		log.Error(err)
	}
}

func Update(db *sql.DB, ap, publishDate string) {
	log.Println("Update violation record with record ...")
	updateDecision := `UPDATE decisions set publish_date = ? WHERE ap = ?`
	statement, err := db.Prepare(updateDecision)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(ap, ap, publishDate)
	if err != nil {
		log.Error(err)
	}
}

func basicGet(db *sql.DB, query string, arg ...interface{}) []*Decision {
	decisions := []*Decision{}
	row, err := db.Query(query, arg...)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		decision := new(Decision)
		err = row.Scan(&decision.AP, &decision.ViolationNumber, &decision.PublishDate)
		if err != nil {
			log.Error(err)
		}
		log.Info(decision)
		decisions = append(decisions, decision)
	}

	return decisions
}

func GetByAll(db *sql.DB) []*Decision {
	return basicGet(db, `SELECT ap, violation_number, publish_date FROM decisions`)
}

func GetByViolationNumber(db *sql.DB, violationNumber string) []*Decision {
	return basicGet(db, `SELECT ap, violation_number, publish_date FROM decisions WHERE violation_number like ? `,
		violationNumber)
}

func GetByDecisionDate(db *sql.DB, decisionDate string) []*Decision {
	return basicGet(db, `SELECT ap, violation_number, publish_date FROM decisions WHERE decision_date like ?`,
		decisionDate)
}
