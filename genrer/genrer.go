package genrer

import (
	"io"
	"os"
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Genrer string

func (g Genrer) String() string { return string(g) }

// lexer guards state of parsing
type lexer struct {
	z	*html.Tokenizer
	token	html.Token
	genres	chan Genrer
}

func (l *lexer) emit(s string) {
	l.genres <- Genrer(s)
}

func Parse(r io.Reader) <-chan Genrer {
	l := &lexer{
		z: html.NewTokenizer(r),
		genres: make(chan Genrer),
	}
	go l.run()
	return l.genres
}

type lexerFunc func(*lexer) lexerFunc

func (l *lexer) run() {
	for f := lexTokenType; f != nil; {
		f = f(l)
	}
	close(l.genres)
}

func lexTokenType(l *lexer) lexerFunc {
	tt := l.z.Next()
	switch tt {
	case html.ErrorToken:
		if err := l.z.Err(); err != io.EOF {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		return nil
	case html.StartTagToken:
		l.token = l.z.Token()
		return lexToken
	}
	return lexTokenType
}

func lexToken(l *lexer) lexerFunc {
	switch l.token.DataAtom {
	case atom.Input:
		for _, attr := range l.token.Attr {
			if attr.Key == "name" && attr.Val == "genres" {
				return lexTokenAttrGenres
			}
		}
	}
	return lexTokenType
}

func lexTokenAttrGenres(l *lexer) lexerFunc {
	for _, attr := range l.token.Attr {
		if attr.Key == "value" {
			l.emit(attr.Val)
		}
	}
	return lexTokenType
}
