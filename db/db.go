package db

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"xkcd/model"

	"github.com/go-sql-driver/mysql"
	mysqlGorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const BATCH_SIZE = 8000

var db *gorm.DB
var err error

func Connect() {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBADDR"),
		DBName:               "xkcd_search",
		AllowNativePasswords: true,
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

func UpdateIncompleteComics(incomplete []model.Comic) {
	//upsert all incomplete comics based on num
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "num"}},
		UpdateAll: true,
	}).Create(&incomplete)
}

func GetLastStoredComicNum() int {
	var comic model.Comic
	db.Last(&comic)
	return comic.Num
}

func GetComics() []model.Comic {
	var comics []model.Comic

	result := db.Find(&comics)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println("GOT ALL COMICS")
	return comics
}

func GetIncomplete() []model.Comic {
	var comics []model.Comic
	result := db.Where("incomplete = ?", true).Find(&comics)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return comics
}

func GetRawWords() []string {
	var title []string
	var transcript []string
	var alt []string
	err := db.Table("comics").Select("title_raw").Find(&title).Error
	if err != nil {
		log.Fatal(err)
	}

	err = db.Table("comics").Select("transcript_raw").Find(&transcript).Error
	if err != nil {
		log.Fatal(err)
	}

	err = db.Table("comics").Select("alt_text_raw").Find(&alt).Error
	if err != nil {
		log.Fatal(err)
	}

	rawString := strings.Join(transcript, " ") + " "
	rawString += strings.Join(title, " ") + " "
	rawString += strings.Join(alt, " ") + " "
	rawWords := strings.Fields(rawString)
	return rawWords
}

func BatchStoreTermFreq(termFreqs []model.TermFreq) {
	termFreqList := make([]model.TermFreqDTO, 0, len(termFreqs))
	for _, termFreq := range termFreqs {
		for term, freq := range termFreq.TermInComicFreq {
			termFreq := model.TermFreqDTO{
				ComicNum: termFreq.ComicNum,
				Term:     term,
				TermsRaw: termFreq.StemToRawMap[term],
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

// Function that updated the termFreq of incomplete comics. Deletes all the old terms and
// insert the new ones
func UpdateTermFreq(termFreqs []model.TermFreq) {
	termFreqList := make([]model.TermFreqDTO, 0, len(termFreqs))
	for _, termFreq := range termFreqs {
		for term, freq := range termFreq.TermInComicFreq {
			termFreq := model.TermFreqDTO{
				ComicNum: termFreq.ComicNum,
				Term:     term,
				TermsRaw: termFreq.StemToRawMap[term],
				Freq:     freq,
			}
			termFreqList = append(termFreqList, termFreq)
		}
		db.Where("comic_num = ?", termFreq.ComicNum).Delete(model.TermFreqDTO{})
	}

	result := db.CreateInBatches(&termFreqList, BATCH_SIZE)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

// Function that returns the frequency of all query terms
func GetTermFreq(queryTerms []string) map[int]model.TermFreq {
	termFreqList := make(map[int]model.TermFreq, 0)

	var termFeqDB []model.TermFreqDTO
	result := db.Where("term in ?", queryTerms).Find(&termFeqDB).Group("comic_num")
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// Loop through the grouped query and store the term frequency struct of the prev comic
	// until the a new comic_num is found (ie all prev comic terms have been stored)
	prevNum := 0
	termFreq := make(map[string]int)
	stemRootMap := make(map[string]string)
	for i, row := range termFeqDB {
		if row.ComicNum != prevNum || i == len(termFeqDB) {
			completedTF := model.TermFreq{
				ComicNum:        prevNum,
				TermInComicFreq: termFreq,
				StemToRawMap:    stemRootMap,
			}
			termFreqList[prevNum] = completedTF
			prevNum = row.ComicNum
			termFreq = make(map[string]int)
			stemRootMap = make(map[string]string)
		}
		termFreq[row.Term] = row.Freq
		stemRootMap[row.Term] = row.TermsRaw
	}
	return termFreqList
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
	result := db.CreateInBatches(&comicFreqList, BATCH_SIZE)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func UpdateComicFreq(newComicFreq model.ComicFreq) {
	comicFreqList := make([]model.ComicFreqDTO, 0, len(newComicFreq.ComicsWithTermFreq))
	for term, freq := range newComicFreq.ComicsWithTermFreq {
		comicFreq := model.ComicFreqDTO{
			Term: term,
			Freq: freq,
		}
		comicFreqList = append(comicFreqList, comicFreq)
	}

	//upsert into comic frequency based on term
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "term"}},
		UpdateAll: true,
	}).CreateInBatches(&comicFreqList, BATCH_SIZE)
}

func GetComicFreq() model.ComicFreq {
	comicFreq := make(map[string]int)

	var comicFreqsDB []model.ComicFreqDTO
	result := db.Find(&comicFreqsDB)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	for _, row := range comicFreqsDB {
		comicFreq[row.Term] = row.Freq
	}
	log.Println("GOT ALL DFS")

	return model.ComicFreq{
		ComicsWithTermFreq: comicFreq,
	}
}
