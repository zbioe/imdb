package title_test

import (
	"io"
	"os"
	"fmt"
	"testing"

	"github.com/iuryfukuda/imdb/title"
)

type parseTest struct {
	parse func(r io.Reader) (<-chan title.Title)
	path  string
	wantCount int
}

var ParserTests = []parseTest{
	musical,
}

var musical = parseTest{
	title.Parse,
	"musical.html",
	100,
}

func TestParse(t *testing.T) {
	for _, test := range ParserTests {
		f := mustLoad(t, test.path)
		var count = 0
		for title := range test.parse(f) {
			count++
			if err := checkTitle(title); err != nil {
				t.Fatal(err)
			}
		}
		if count == 0 {
			t.Fatalf("not found")
		}
		if test.wantCount != count {
			t.Fatalf("count: want %d; got %d", test.wantCount, count)
		}
	}
}

func checkTitle(t title.Title) error {
	if len(t.Genres) == 0 {
		return fmt.Errorf("can't find genres in title: %#v", t)
	}
	if t.Text == "" {
		return fmt.Errorf("can't find text in title: %#v", t)
	}
	if t.Year == "" {
		return fmt.Errorf("can't find year in title: %#v", t)
	}
	r := t.Rating
	if r.Value == 0.0 {
		return fmt.Errorf("can't find value in title.rating: %#v", r)
	}
	if r.Best == 0.0 {
		return fmt.Errorf("can't find best in title.rating: %#v", r)
	}
	if r.Count == 0 {
		return fmt.Errorf("can't find count in title.rating: %#v", r)
	}
	if r.Position == 0 {
		return fmt.Errorf("can't find position in title.rating: %#v", r)
	}
	return nil
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