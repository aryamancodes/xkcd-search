package main

import (
	"xkcd/api"
	"xkcd/db"
	"xkcd/index"
)

func main() {
	//populateDB()
	api.Serve()
}

func populateDB() {
	db.Connect()
	comics := index.FetchAllComics()
	tf := index.ComputeAllTermFreq(comics)
	index.ComputeAllComicFreq(comics, tf)
}
