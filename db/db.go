/*
Copyright © 2025 Thomas McKinstry thomas.g.mckinstry@protonmail.com
*/

package db

import (
	"database/sql"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDB() *sql.DB {
	var err error

	if db != nil {
		//utils.DebugLog("Database already exists", nil)
		return db
	}

	db, err = sql.Open("sqlite3", utils.Config.Database)
	utils.DebugLog("Created new database instance", nil)
	if err != nil {
		log.Fatal("Unable to open database:", err)

		os.Exit(1)
	}
	init_db(db)

	return db
}
