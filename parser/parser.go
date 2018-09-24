package parser

import (
	"log"
	"io"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseGenres(r io.Reader) ([]string, error) {
	var (
		genres []string
		z = html.NewTokenizer(r)
	)
Loop:
	for {
		tokenType := z.Next()
		token := z.Token()
		switch tokenType {
		case html.ErrorToken:
			switch err := z.Err(); err {
			case io.EOF:
				break Loop
			}
		case html.StartTagToken:
			switch token.DataAtom {
			case atom.Input:
				for _, attr := range token.Attr {
					if attr == "name" {
						log.Printf("here: %v", attr)
					}
				}
				continue Loop
			}
		}
	}
	return genres, nil
}
