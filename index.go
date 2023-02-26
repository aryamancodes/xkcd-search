package main

import (
	"io/ioutil"
	"log"
	"net/http"
 )


func main() {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
	   log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
	
	sb := string(body)
	log.Printf(sb)
 }