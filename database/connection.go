package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the *gorm.DB instance used by the application for persistent storage.
var DB *gorm.DB

// ConnectToSQLiteDatabase opens a SQLite connection to the given file and
// returns the created *gorm.DB object.
func ConnectToSQLiteDatabase(fileName string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(fileName), &gorm.Config{
		TranslateError: true,
	})
}

// ConnectToTestDatabase opens an in-memory SQLite database connection and
// returns the created *gorm.DB object.
func ConnectToTestDatabase() (*gorm.DB, error) {
	return ConnectToSQLiteDatabase(":memory:")
}

// InitDBOrPanic connects to the database and sets the DB package variable to
// a pointer to the created gorm.DB. If an error occurs, the application will
// print the error and panic.
func InitDBOrPanic() {
	db, err := ConnectToSQLiteDatabase("data.sqlite.db")
	if err != nil {
		log.Panic(err)
	}
	DB = db
}
