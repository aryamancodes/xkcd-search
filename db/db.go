package db

import (
	"database/sql"
	"log"
	"os"
	"xkcd/model"

	"github.com/go-sql-driver/mysql"
	mysqlGorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const BATCH_SIZE = 6000

var db *gorm.DB
var err error

func Connect() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "xkcd_search",
	}
	sqlDB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(mysqlGorm.New(mysqlGorm.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected @ port 3306!")
}

func BatchStoreComics(comics []model.Comic) {
	result := db.Create(&comics)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func BatchStoreTermFreq(termFreqs []model.TermFreq) {
	termFreqList := make([]model.TermFreqDTO, 0, len(termFreqs))
	for _, termFreq := range termFreqs {
		for term, freq := range termFreq.TermInComicFreq {
			termFreq := model.TermFreqDTO{
				ComicNum: termFreq.Comic.Num,
				Term:     term,
				Freq:     freq,
			}
			termFreqList = append(termFreqList, termFreq)
		}
	}
	result := db.CreateInBatches(&termFreqList, BATCH_SIZE)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func BatchStoreComicFreq(comicFreqs model.ComicFreq) {
	comicFreqList := make([]model.ComicFreqDTO, 0, len(comicFreqs.ComicsWithTermFreq))
	for term, freq := range comicFreqs.ComicsWithTermFreq {
		comicFreq := model.ComicFreqDTO{
			Term: term,
			Freq: freq,
		}
		comicFreqList = append(comicFreqList, comicFreq)
	}
	db.Logger.LogMode(logger.Info)
	result := db.CreateInBatches(&comicFreqList, BATCH_SIZE)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}
