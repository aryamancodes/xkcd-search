package main

import (
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

func init() {
	api.Serve()
}

func main() {
	//populateDB()
	//api.Serve()
	lambda.Start(api.AWSHandler)
}

func populateDB() {
	db.Connect()
	comics := index.FetchAllComics()
	tf := index.ComputeAllTermFreq(comics)
	index.ComputeAllComicFreq(comics, tf)
}
