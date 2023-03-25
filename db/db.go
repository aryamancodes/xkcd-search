package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Connect() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "xkcd_search",
	}
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(15)
	log.Println("Connected @ port 3306!")
}

// TODO: Batch insert!!
func StoreComics(num int, text string) {
	_, err := db.Exec("INSERT INTO comics (num, transcript) VALUES (?,?)", num, text)
	if err != nil {
		log.Fatal(err)
	}
}

func StoreTermFreq(num int, term string, freq int, total int) {
	_, err := db.Exec("INSERT INTO term_frequency (comic_num, term, termFreq, totalTerms) VALUES (?,?,?,?)", num, term, freq, total)
	if err != nil {
		log.Fatal(err)
	}
}

func StoreComicFreq(term string, freq int, total int) {
	_, err := db.Exec("INSERT INTO comic_frequency (term, comicFreq, totalComics) VALUES (?,?,?)", term, freq, total)
	if err != nil {
		log.Fatal(err)
	}
}
