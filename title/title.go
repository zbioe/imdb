package title

import (
	"io"
	"log"
	"strings"
	"unicode"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Title struct {
	Name	string
	Episode	string	`json:",omitempty"`
	Year	string
	Genres	[]string
	Rating	Rating
}

type Rating struct {
	Value	float64	`json:",omitempty"`
	Best	float64	`json:",omitempty"`
	Count	int	`json:",omitempty"`
	Position	int	`json:",omitempty"`
}

type Result struct {
	Titles	<-chan Title
	Error	error
}

func Parse(r io.Reader) Result {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return Result{Error:err}
	}
	titles := make(chan Title)
	go runParse(doc, titles)
	return Result{Titles:titles}
}

func runParse(doc *goquery.Document, titles chan<- Title) {
	doc.Find(".mode-advanced").Each(func(i int, s *goquery.Selection) {
		titles <- parseTitle(s)
	})
	close(titles)
}

func parseTitle(s *goquery.Selection) Title {
	var title Title
	title.Name =  parseTitleName(s)
	title.Episode = parseTitleEpisode(s)
	title.Year = parseTitleYear(s)
	title.Genres = parseTitleGenres(s)
	title.Rating = parseTitleRating(s)
	return title
}


func parseTitleName(s *goquery.Selection) string {
	var title string
	div := s.Find(".lister-item-header")
	a := div.Find(`a[href#=(/title/)]`)
	title = strings.Trim(a.First().Text(), " ")
	return trimLetter(title)
}

func parseTitleEpisode(s *goquery.Selection) string {
	var episode string
	div := s.Find(".lister-item-header")
	a := div.Find("a[href#=(/title/)]")
	if a.Size() > 1 {
		episode = a.Last().Text()
	}
	return trimLetter(episode)
}

func parseTitleYear(s *goquery.Selection) string {
	var year string
	div := s.Find(".lister-item-year")
	year = div.First().Text()
	return trimSpace(year)
}

func parseTitleGenres(s *goquery.Selection) []string {
	var genres []string
	div := s.Find(".genre")
	for _, g := range strings.SplitN(div.Text(), ",", -1) {
		genrer := g
		genrer = trimLetter(genrer)
		genrer = strings.ToLower(genrer)
		genres = append(genres, genrer)
	}
	return genres
}

func parseTitleRating(s *goquery.Selection) Rating {
	var rating Rating
	rating.Value = parseTitleRatingValue(s)
	rating.Best = parseTitleRatingBest(s)
	rating.Count = parseTitleRatingCount(s)
	rating.Position = parseTitleRatingPosition(s)
	return rating
}

func parseTitleRatingValue(s *goquery.Selection) float64 {
	var value float64
	meta := s.Find("[itemprop=ratingValue]")
	rawvalue, ok := meta.Attr("content")
	if ok {
		value = toFloat64(rawvalue)
	}
	return value
}

func parseTitleRatingBest(s *goquery.Selection) float64 {
	var best float64
	meta := s.Find("[itemprop=bestRating]")
	rawbest, ok := meta.Attr("content")
	if ok {
		best = toFloat64(rawbest)
	}
	return best
}

func parseTitleRatingCount(s *goquery.Selection) int {
	var count int
	meta := s.Find("[itemprop=ratingCount]")
	rawcount, ok := meta.Attr("content")
	if ok {
		count = toInt(rawcount)
	}
	return count
}

func parseTitleRatingPosition(s *goquery.Selection) int {
	var position int
	div := s.Find(".lister-item-header")
	raw := div.Find(".lister-item-index").Text()
	position = toInt(strings.Replace(raw, ",", "", 1))
	return position
}

func toInt(s string) int {
	rawint := trimNumber(s)
	i, err := strconv.Atoi(rawint)
	if err != nil {
		log.Printf("Can't converto to int: %s\n", err)
		return 0
	}
	return i
}

func toFloat64(s string) float64 {
	rawfloat := trimNumber(s)
	f, err := strconv.ParseFloat(rawfloat, 64)
	if err != nil {
		log.Printf("Can't converto to float: %s\n", err)
		return 0
	}
	return f
}

func trimNumber(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsNumber(r)
	})
}

func trimLetter(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

func trimSpace(s string) string {
	return strings.Trim(s, "\n \t")
}