package db

import (
	"database/sql"
	"log"

	"github.com/thomasmckinstry/ouevre/utils"
)

func init_db(db *sql.DB) {
	var err error

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS works (
	date_added date,
	title varchar(80) NOT NULL,
	media_type text NOT NULL, -- Use golang JSON Marshal/Unmarshal to make the arrays into strings
	work_status integer NOT NULL DEFAULT 0,
	tags text, -- Marshalled array
	year_released integer NOT NULL,
	work_id integer PRIMARY KEY -- Need work_id to specify between adaptations with the same title (ex. Running Man movie v. Running Man book)
);
	`)
	if err != nil {
		log.Fatal("Unable to create works table in database: ", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reviews (
	date_added date,
	review_text text NOT NULL,
	review_id integer PRIMARY KEY,
	work_id integer REFERENCES works
);
	`)
	if err != nil {
		log.Fatal("Unable to create reviews table in database: ", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS notes (
	date_added date,
	note_text text NOT NULL,
	note_id integer PRIMARY KEY,
	work_id integer REFERENCES works
);

	`)
	if err != nil {
		log.Fatal("Unable to create notes table in database: ", err)
	}

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS status_table (
				id integer PRIMARY KEY,
				status_name varchar(15) UNIQUE 
			);
		`)
	utils.CheckError("Unable to create status_table in database: ", err)

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS media_type_table (
				id integer PRIMARY KEY,
				type_name varchar(25) UNIQUE
			);
		`)
	utils.CheckError("Unable to create media_type_table in database: ", err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tags_table (
	tag_name varchar(25) NOT NULL PRIMARY KEY
);
	`)
	utils.CheckError("Unable to create table in database: ", err)
}
