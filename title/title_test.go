package title_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/iuryfukuda/imdb/title"
)

type parseTest struct {
	raw     string
	want    []title.Title
}

func runParseTest(t parseTest) error {
		result := title.Parse(strings.NewReader(t.raw))
		if result.Error != nil {
			return fmt.Errorf("unexpected err: %s", result.Error)
		}

		for i, want := range t.want {
			got, ok := <-result.Titles
			if !ok || !reflect.DeepEqual(got, want) {
				return fmt.Errorf("[%d]: got [%#v], want [%#v]", i, got, want)
			}
		}

		title, ok := <-result.Titles
		if ok {
			return fmt.Errorf("channel still open and returns: %#v", title)
		}

		return nil
}

func BenchmarkParse(b *testing.B) {
	for _, test := range parseTests {
		if err := runParseTest(test); err != nil {
			b.Fatal(err)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		if err := runParseTest(test); err != nil {
			t.Fatal(err)
		}
	}
}

var parseTests = []parseTest{
	parseTest{
		raw: `
<div class="lister-item mode-advanced">
    <div class="lister-item-content">
		<h3 class="lister-item-header">
        	<span class="lister-item-index unbold text-primary">
        		3,819.
        	</span>
    		<a href="/title/tt0085959/?ref_=adv_li_tt">
    			Monty Python - O Sentido da Vida
    		</a>
			<span class="lister-item-year text-muted unbold">
				(1983)
			</span>
		</h3>
    	<p class="text-muted ">
            <span class="certificate">18</span>
			<span class="runtime">107 min</span>
            <span class="genre">Comedy, Musical</span>
    	</p>
    </div>
	<div class="inline-block ratings-user-rating">
		<div>
			<meta itemprop="ratingValue" content="7.6" />
			<meta itemprop="bestRating" content="10" />
			<meta itemprop="ratingCount" content="99891" />
		</div>
	</div>
</div>
<div class="lister-item mode-advanced">
	<div class="lister-item-content">
		<h3 class="lister-item-header">
			<span class="lister-item-index unbold text-primary">
				3,801.
			</span>
			<a href="/title/tt4902964/?ref_=adv_li_tt">
				Enrolados Outra Vez
			</a>
			<span class="lister-item-year text-muted unbold">
				(2017– )
			</span>
        	<br />
			<small class="text-primary unbold">Episode:</small>
			<a href="/title/tt6581940/?ref_=adv_li_tt">
				What the Hair?!
			</a>
			<span class="lister-item-year text-muted unbold">
				(2017)
			</span>
		</h3>
		<p class="text-muted ">
			<span class="runtime">22 min</span>
            <span class="genre">
				Animation, Adventure, Comedy
			</span>
		</p>
	</div>
 	<div class="inline-block ratings-user-rating">
    <div class="starBarWidget" id="sb_tt6581940">
		<div>
			<meta itemprop="ratingValue" content="7.7" />
			<meta itemprop="bestRating" content="10" />
			<meta itemprop="ratingCount" content="72" />
		</div>
	</div>
</div>
		`,
		want: []title.Title{
			title.Title{
				Name: "Monty Python - O Sentido da Vida",
				Year: "(1983)",
				Genres: []string{"comedy", "musical"},
				Rating: title.Rating{
					Value:7.6, Best:10, Count:99891, Position:3819,
				},
			},
			title.Title{
				Name: "Enrolados Outra Vez",
				Episode: "What the Hair",
				Year: "(2017– )",
				Genres: []string{"animation", "adventure", "comedy"},
				Rating: title.Rating{
					Value:7.7, Best:10, Count:72, Position:3801,
				},
			},
		},
	},
	parseTest{
		raw: "",
		want: []title.Title{},
	},
}
