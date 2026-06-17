/*
Copyright © 2025 Thomas McKinstry thomas.g.mckinstry@protonmail.com
*/

package db

import (
	"database/sql"
	"github.com/thomasmckinstry/ouevre/utils"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDB() *sql.DB {
	var err error

	if db != nil {
		//utils.DebugLog("Database already exists", nil)
		return db
	}

	db, err = sql.Open("sqlite3", utils.DirectoryPath+utils.Config.Database)
	utils.DebugLog("Created new database instance", nil)
	utils.CheckError("Unable to open database:", err)
	init_db(db)

	return db
}

func SetStatuses(statuses []string) {
	db = GetDB()
	query, err := db.Prepare(`INSERT OR IGNORE INTO status_table (status_name) VALUES (?);`)
	utils.CheckError("Failed to prepare query for insert statuses: ", err)
	for _, status := range statuses {
		_, err = query.Exec(status)
		utils.CheckError("Failed to insert status: ", err)
	}
	query.Close()
}

func SetMediums(mediums []string) {
	db = GetDB()
	query, err := db.Prepare(`INSERT OR IGNORE INTO media_type_table (type_name) VALUES (?);`)
	utils.CheckError("Failed to prepare query for insert mediums: ", err)
	defer query.Close()
	for _, medium := range mediums {
		_, err = query.Exec(medium)
		utils.CheckError("Failed to insert medium: ", err)
	}
	query.Close()
}
