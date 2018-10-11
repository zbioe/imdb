package genrer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/iuryfukuda/imdb/genrer"
)

type parseTest struct {
	raw  string
	want  []string
}

var parseTests = []parseTest{
	parseTest{
		raw: `  <table>
            <tbody>
            <tr>
              <td><input id="genres-1" type="checkbox" name="genres" value="action"> <label for="genres-1">Action</label></td>
              <td><input id="genres-2" type="checkbox" name="genres" value="adventure"> <label for="genres-2">Adventure</label></td>
              <td><input id="genres-3" type="checkbox" name="genres" value="animation"> <label for="genres-3">Animation</label></td>
              <td><input id="genres-4" type="checkbox" name="genres" value="biography"> <label for="genres-4">Biography</label></td>
            </tr>
            <tr>
		`,
		want: []string{
			"action", "adventure", "animation", "biography",
		},
	},
	parseTest{
		raw: "",
		want: []string{},
	},
}

func runParseTest(t parseTest) error {
		genres := genrer.Parse(strings.NewReader(t.raw))

		for i, want := range t.want {
			got, ok := <-genres
			if !ok || got.String() != want {
				return fmt.Errorf("[%d]: got [%#v], want [%#v]", i, got, want)
			}
		}

		genre, ok := <-genres
		if ok {
			return fmt.Errorf("channel still open and returns: %#v", genre)
		}

		return nil
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		if err := runParseTest(test); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for _, test := range parseTests {
		if err := runParseTest(test); err != nil {
			b.Fatal(err)
		}
	}
}