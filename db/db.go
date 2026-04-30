/*
Copyright © 2025 Thomas McKinstry thomas.g.mckinstry@protonmail.com
*/

package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var configs struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

var db *sql.DB

func init() {
	content, _ := os.ReadFile("config.json") // TODO: Add in the util for checking if a sql query has worked
	//c.Check(err, "ERROR: Failed to read config file.")

	_ = json.Unmarshal(content, &configs)
	//c.Check(err, "ERROR: Failed to unmarshal config.")
}

func GetWorkID(db *sql.DB) int {
	var id int
	rows, err := db.Query(`SELECT work_id FROM works ORDER BY work_id DESC LIMIT 1;`)
	if err != nil {
		log.Fatal("Unable to query database for id: ", err)
	}
	defer rows.Close()
	rows.Next()
	if rows == nil {
		return 0
	}
	err = rows.Scan(&id)
	if err != nil {
		return 0 // TODO: This is a terrible hack for when there are no existing rows. Need to do something about it
		//log.Fatal("Unable to scan rows for id: ", err)
	}

	if len(os.Getenv("DEBUG")) > 0 {
		log.Println("Got id: ", id+1)
	}
	return id + 1
}

func GetDB() *sql.DB {
	var err error

	if db != nil {
		return db
	}

	db, err := sql.Open("sqlite3", "./media.db")
	if err != nil {
		log.Fatal("Unable to open database:", err)

		os.Exit(1)
	}
	init_db(db)

	return db
}
