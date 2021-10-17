package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func init() {

	log.Info("Creating police.db...")
	file, err := os.Create("police.db") // Create SQLite file
	if err != nil {
		log.Error(err)
	}
	file.Close()
	log.Info("police.db created")

	db, err = sql.Open("sqlite3", "./police.db") // Open the created SQLite File
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer db.Close() // Defer Closing the database
	createTables()   // Create Database Tables
}

func createTables() {
	createTableViolations()
}

func createTableViolations() {
	createViolationTableSQL := `CREATE TABLE violations (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" VARCHAR,
		"name" VARCHAR		
	  );` // SQL Statement for Create Table

	log.Info("Create violations table...")
	statement, err := db.Prepare(createViolationTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Error(err)
	}
	statement.Exec() // Execute SQL Statements
	log.Info("violations table created")
}

func GetDB() *sql.DB {
	return db
}
