package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"xkcd/db"
 )

func main() {
	resp, err := http.Get("https://xkcd.com/2741/info.0.json")
	if err != nil {
	   log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
	
	db.New()

	body_string := string(body)
	log.Printf(body_string)
 }