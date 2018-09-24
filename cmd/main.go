package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

const SearchURL = "https://www.imdb.com/search/title"

var genrer = "all"
var limitiItems = 500
var adult = "include"
var sort = "user_rating,desc"
var itemsPerPage = 100

func main() {
	setup()
	resp, err := http.Get(SearchURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("imdb: ")
}
