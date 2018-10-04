package genrer_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/iuryfukuda/imdb/genrer"
)

type parseTest struct {
	parse func(r io.Reader) (<-chan genrer.Genrer)
	path  string
	want  []string
}

var ParserTests = []parseTest{
	genres,
}

var genres = parseTest{
	parse: genrer.Parse,
	path:  "genres.html",
	want: []string{
		"action", "adventure", "animation", "biography",
		"comedy", "crime", "documentary", "drama",
		"family", "fantasy", "film_noir", "game_show",
		"history", "horror", "music", "musical",
		"mystery", "news", "reality_tv", "romance",
		"sci_fi", "sport", "talk_show", "thriller",
		"war", "western",
	},
}

func TestParse(t *testing.T) {
	for _, test := range ParserTests {
		f := mustLoad(t, test.path)
		var got []string
		for g := range test.parse(f) {
			got = append(got, g.String())
		}
		f.Close()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got: %s; want: %s", got, test.want)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for _, test := range ParserTests {
		for n := 0; n < b.N; n++ {
			var got []string
			f := mustLoadB(b, test.path)
			for g := range test.parse(f) {
				got = append(got, g.String())
			}
			if !reflect.DeepEqual(got, test.want) {
				b.Errorf("got: %s; want: %s", got, test.want)
			}
			f.Close()
		}
	}
}

func load(path string) (*os.File, error) {
	const prefixPath = "./testdata/"
	return os.Open(prefixPath + path)
}

func mustLoad(t *testing.T, path string) *os.File {
	f, err := load(path)
	if err != nil {
		t.Fatalf("Can't load file %s: %s\n", path, err)
	}
	return f
}

func mustLoadB(b *testing.B, path string) *os.File {
	f, err := load(path)
	if err != nil {
		b.Fatalf("Can't load file %s: %s\n", path, err)
	}
	return f
}
