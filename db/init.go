package db

import (
	"database/sql"
	"log"
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
	work_id integer PRIMARY KEY
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reviews (
	date_added date,
	review_text text NOT NULL,
	review_score integer NOT NULL check (review_score <= 10),
	review_id integer PRIMARY KEY,
	work_id integer REFERENCES works
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
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
		log.Fatal("Unable to create table in database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS status_table (
	id integer PRIMARY KEY,
	status_name varchar(15) NOT NULL
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS media_type_table (
	id integer PRIMARY KEY,
	type_name varchar(25) NOT NULL
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS genre_table (
	genre_name varchar(25) NOT NULL PRIMARY KEY
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS theme_table (
	theme_name varchar(25) NOT NULL UNIQUE PRIMARY KEY
);
	`)
	if err != nil {
		log.Fatal("Unable to create table in database:", err)
	}
}
