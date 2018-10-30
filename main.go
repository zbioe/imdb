package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"math"
	"sync"
	"strings"
	"net/url"
	"net/http"
	"encoding/json"

	"github.com/iuryfukuda/imdb/title"
	"github.com/iuryfukuda/imdb/genrer"
)

const resultPATH = "./results"

var (
	SearchURL = "https://www.imdb.com/search/title"
	limit = flag.Int("limit", 500, "limit per genrer")
	adult = flag.Bool("adult", true, "incluse adult results")
	debug = flag.Bool("debug", false, "verbose debug mode")
	sort = flag.String("sort", "user_rating,desc", "sorted by")
	itemsPerReq = flag.Int("count", 50, "items returned per request")
)

func main() {
	flag.Parse()
	setup()

	resp, err := http.Get(SearchURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err := os.Stat(resultPATH); os.IsNotExist(err) {
		err := os.Mkdir(resultPATH, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	var query = url.Values{}
	for _, s := range strings.SplitN(*sort, ",", -1) {
		query.Add("sort", s)
	}
	if *adult {
		query.Add("adult", "include")
	}

	var wg = new(sync.WaitGroup)
	var rawquery = query.Encode()
	for g := range genrer.Parse(resp.Body) {
		var rq = fmt.Sprintf("%s&genres=%s", rawquery, g)
		if *debug {
			log.Printf("start collect %s", g)
		}
		wg.Add(1)
		go collectTitles(wg, g, *debug, rq, *limit, *itemsPerReq)
	}
	wg.Wait()
}


func collectTitles(
	wg *sync.WaitGroup,
	g genrer.Genrer,
	debug bool,
	rawquery string,
	limit, itemsPerReq int,
) {
	defer wg.Done()
	var rq = fmt.Sprintf("%s&count=%d", rawquery, itemsPerReq)
	var sum	int
	var npage = calculatePages(limit, itemsPerReq)
	var filepath = fmt.Sprintf("%s/%s.jsonl", resultPATH, g)

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	for p := 1; p <= npage; p++ {
		var rq = fmt.Sprintf("%s&page=%d", rq, p)
		var rawurl = fmt.Sprintf("%s?%s", SearchURL, rq)

		if debug {
			log.Printf("send request to %s", rawurl)
		}
		resp, err := http.Get(rawurl)
		if err != nil {
			log.Fatal(err)
		}

		if debug {
			log.Printf("start process titles of %s", rawurl)
		}
		result := title.Parse(resp.Body)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		for t := range result.Titles {
			sum++
			if debug {
				log.Printf("%s: process %dº title of %dº page", g, sum, p)
			}
			if err := encoder.Encode(t); err != nil {
				log.Fatal(err)
			}
			if sum == limit {
				break
			}
		}

		if debug {
			log.Printf("finish process titles of %s", rawurl)
		}
		resp.Body.Close()
	}
	if debug {
		log.Printf("finish collect %s", g)
	}
}

func calculatePages(limit, itemsPerReq int) int {
	var npage int
	divided := float64(limit) / float64(itemsPerReq)
	if truncated := math.Trunc(divided); truncated == divided {
		npage = int(truncated)
	} else {
		npage = int(truncated) + 1
	}
	return npage
}

func setup() {
	log.SetFlags(0)
	log.SetPrefix("imdb: ")
}
