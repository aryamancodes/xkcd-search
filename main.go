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
var isUpdate = false

func init() {
	runLocal = os.Getenv("AWS_LAMBDA_RUNTIME_API") == ""

	//if any command line arg is passed, we're doing a scheduled updated
	isUpdate = len(os.Args) > 1
	if !isUpdate {
		api.Serve(runLocal)
	}
}

func main() {
	//populateDB()
	if isUpdate {
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
	var updatedComicFreq model.ComicFreq

	//handle incomplete comics that might have now been updated
	oldComics := db.GetIncomplete()
	if len(oldComics) > 0 {
		currComics := index.FetchIncompleteComics(oldComics)
		revisedTf := index.ComputeIncompleteTermFreq(oldComics, currComics)
		updatedComicFreq = index.ComputeNewComicFreq(revisedTf, db.GetComicFreq())
	} else {
		updatedComicFreq = db.GetComicFreq()
	}

	//handle new comics that have been added since last update
	newComics := index.FetchNewComics()
	if newComics != nil {
		newTf := index.ComputeAllTermFreq(newComics)
		index.ComputeNewComicFreq(newTf, updatedComicFreq)
	}
}
