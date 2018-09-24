package parser_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/iuryfukuda/imdb/parser"
)

type parserTest struct {
	parse func(r io.Reader) ([]string, error)
	path  string
	want  []string
	err   error
}

var allParserTests = []parserTest{
	genres,
}

var genres = parserTest{
	parse: parser.ParseGenres,
	path:  "title.html",
	want: []string{
		"action", "adventure", "animation", "biography",
		"comedy", "crime", "documentary", "drama",
		"family", "fantasy", "film-noir", "game-show",
		"history", "horror", "music", "musical",
		"mystery", "news", "reality-tv", "romance",
		"sci-fi", "sport", "talk-show", "thriller",
		"war", "western",
	},
	err: nil,
}

func TestParser(t *testing.T) {
	for _, test := range allParserTests {
		f := mustLoad(t, test.path)
		got, err := test.parse(f)
		f.Close()
		if err != test.err {
			t.Errorf("got: error == %s; want: error == %s", err, test.err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got: %s; want: %s", got, test.want)
		}
	}
}

func mustLoad(t *testing.T, path string) *os.File {
	const prefixPath = "./testdata/"
	f, err := os.Open(prefixPath + path)
	if err != nil {
		t.Fatalf("Can't load file %s: %s\n", path, err)
	}
	return f
}
