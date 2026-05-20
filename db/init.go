package db

import (
	"database/sql"
	"log"

	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
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
	work_id integer UNIQUE -- Need work_id to specify between adaptations with the same title (ex. Running Man movie v. Running Man book)
	PRIMARY KEY (title, media_type)
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

	row, err := db.Query(`
		SELECT 1 FROM sqlite_master WHERE type='table' AND name='status_table';
	`)
	utils.CheckError("Unable to determine if status_table exists: ", err)
	if !row.Next() {
		_, err = db.Exec(`
			CREATE TABLE status_table (
				id integer PRIMARY KEY,
				status_name varchar(15) NOT NULL
			);
		`)
		utils.CheckError("Unable to create status_table in database: ", err)

		_, err = db.Exec(`
			INSERT INTO status_table (id, status_name)
			VALUES
				(0, 'Pending'),
				(1, 'Started'),
				(2, 'Hiatus'),
				(3, 'Completed'),
				(4, 'Dropped');
			`)
		utils.CheckError("Unable to insert to status_table in database: ", err)
	}
	err = row.Close()
	utils.CheckError("Failed to close insert to status_table: ", err)

	// TODO: The prefill on this table should probably removed, I should set up some config options or something so people can customize
	// Could probably just keep an array in memory from a config file and index into it to convert
	row, err = db.Query(`SELECT 1 from sqlite_master WHERE type='table' AND name='media_type_table'`)
	utils.CheckError("Unable to determine if media_type_table exists: ", err)
	if !row.Next() {
		_, err = db.Exec(`
			CREATE TABLE media_type_table (
				id integer PRIMARY KEY,
				type_name varchar(25) NOT NULL
			);
		`)
		utils.CheckError("Unable to create media_type_table in database: ", err)

		_, err = db.Exec(`
			INSERT INTO media_type_table (id, type_name)
			VALUES
				(0, 'Anime'),
				(1, 'Manga'),
				(2, 'Movie'),
				(3, 'Book'),
				(4, 'Comic'),
				(5, 'Show'),
				(6, 'Animated'),
				(7, 'Live Action');
		`)
		utils.CheckError("Unable to insert to media_type_table in database: ", err)
	}
	err = row.Close()
	utils.CheckError("Failed to close insert to media_type_table: ", err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tags_table (
	tag_name varchar(25) NOT NULL PRIMARY KEY
);
	`)
	utils.CheckError("Unable to create table in database: ", err)
}
