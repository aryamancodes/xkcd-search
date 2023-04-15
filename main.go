package main

import (
	"os"
	"xkcd/api"
	"xkcd/db"
	"xkcd/index"
	"xkcd/model"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/sajari/fuzzy"
)

var comicFreq model.ComicFreq
var language *fuzzy.Model
var ginLambda *ginadapter.GinLambda
var runLocal = false
var isSchedule = false

func init() {
	runLocal = os.Getenv("AWS_LAMBDA_RUNTIME_API") == ""
	isSchedule = os.Getenv("SCHEDULE") != ""
	if !isSchedule {
		api.Serve(runLocal)
	}
}

func main() {
	//populateDB()
	if isSchedule {
		updateData()
	} else if !runLocal {
		lambda.Start(api.AWSHandler)
	}
}

func populateDB() {
	db.Connect()
	comics := index.FetchAllComics()
	tf := index.ComputeAllTermFreq(comics)
	index.ComputeAllComicFreq(comics, tf)
}

func updateData() {
	db.Connect()
	//handle incomplete comics that might have now been updated
	incompleteComics := db.GetIncomplete()

	//handle new comics that have been added since last update
	newComics := index.FetchNewComics()
	newTf := index.ComputeAllTermFreq(newComics)

	index.ComputeNewComicFreq(newTf, db.GetComicFreq())
}
