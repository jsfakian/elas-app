package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func init() {

	_, err := os.Stat("./police.db")
	creatingDB := os.IsNotExist(err)
	if creatingDB {
		log.Info("Creating police.db...")
		file, err := os.Create("police.db") // Create SQLite file
		if err != nil {
			log.Error(err)
		}
		file.Close()
		log.Info("police.db created")
	}

	db, err = sql.Open("sqlite3", "./police.db") // Open the created SQLite File
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if creatingDB {
		createTables() // Create Database Tables
	}
}

func createTables() {
	createTableViolations()
	createTableOfficer()
	createTableDecisions()
	createTableObjections()
}

func createTable(tableDescription string) {
	log.Info("Creating table...")
	statement, err := db.Prepare(tableDescription) // Prepare SQL Statement
	if err != nil {
		log.Error(err)
	}
	statement.Exec() // Execute SQL Statements
	log.Info("Table created")
}

func createTableOfficer() {
	createOfficerTableSQL := `CREATE TABLE officer (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"first_name" VARCHAR,
		"last_name" VARCHAR,
		"rank" VARCHAR,
		"commander", INTEGER		
	  );` // SQL Statement for Create Table

	createTable(createOfficerTableSQL)
}

func createTableViolations() {
	createViolationTableSQL := `CREATE TABLE violations (
		"ap" VARCHAR NOT NULL PRIMARY KEY,
		"at" VARCHAR,
		"violation_number" INTEGER,
		"registration_number" VARCHAR,
		"first_name_owner" VARCHAR,
		"middle_name_owner" VARCHAR,
		"last_name_owner" VARCHAR,
		"address_owner" VARCHAR,
		"document_type" INTEGER
	  );` // SQL Statement for Create Table

	createTable(createViolationTableSQL)
}

func createTableDecisions() {
	createDecisionTableSQL := `CREATE TABLE decisions (
		"ap" VARCHAR NOT NULL PRIMARY KEY,
		"violation_number" INTEGER,
		"publish_date" VARCHAR
	  );` // SQL Statement for Create Table

	createTable(createDecisionTableSQL)
}

func createTableObjections() {
	createObjectionTableSQL := `CREATE TABLE objections (
		"ap" VARCHAR NOT NULL PRIMARY KEY,
		"violation_number" INTEGER,
		"first_name_driver" VARCHAR,
		"last_name_driver" VARCHAR,
		"middle_name_driver" VARCHAR,
		"objection_date" VARCHAR,
		"publish_date" VARCHAR
	  );` // SQL Statement for Create Table

	createTable(createObjectionTableSQL)
}

func GetDB() *sql.DB {
	return db
}
