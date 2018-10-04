package main

import (
	"log"
	"net/http"
	"net/url"
	"io/ioutil"

	"github.com/iuryfukuda/imdb/parser"
)

const SearchURL = "https://www.imdb.com/search/title"

// var genrers = []string{"all"}
var limitiItems = 500
var adult = "include"
var sort = "user_rating,desc"
var itemsPerPage = 100
var baseURL, _ = url.Parse(SearchURL)

func main() {
	setup()
	genres, err := getGenres()
	if err != nil {
		fail(err)
	}
	for _, g := range genres {
		titles, err := getTitles(g)
		if err != nil {
			fail(err)
		}
		err = toFile(g, titles)
		if err != nil {
			fail(err)
		}
	}
	
}

func toFile(genre, titles string) error {
	log.Print(genre, titles)
	return nil
}

func getGenres() ([]string, error) {
	resp, err := http.Get(SearchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parser.ParseGenres(resp.Body)
}

func getTitles(genrer string) (string, error){
	v := url.Values{}
	v.Add("genres", genrer)
	v.Add("adult", "include")
	v.Add("sort", "user_rating")
	v.Add("sort", "desc")
	v.Add("count", "100")
	u := baseURL
	u.RawQuery = v.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func fail(err error) {
	log.Fatal(err)
}

func setup() {
	log.SetFlags(0)
	log.SetPrefix("imdb: ")
}
