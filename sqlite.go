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
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"first_name" VARCHAR,
		"last_name" VARCHAR,
		"rank" VARCHAR,
		"commander", INTEGER		
	  );` // SQL Statement for Create Table

	createTable(createOfficerTableSQL)
}

func createTableViolations() {
	createViolationTableSQL := `CREATE TABLE violations (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" VARCHAR,
		"name" VARCHAR		
	  );` // SQL Statement for Create Table

	createTable(createViolationTableSQL)
}

func GetDB() *sql.DB {
	return db
}
