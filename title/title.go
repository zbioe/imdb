package title

import (
	"io"
	"os"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Title struct {
	Genres []string
	Text	string
	Year	string
	Rating	Rating
}

type Rating struct {
	Value	float64
	Best	float64	
	Count	int
	Position	int
}

// lexer guards state of parsing
type lexer struct {
	depth	int
	emitDepth	int
	// guards html types of tokenizer
	z	*html.Tokenizer
	t	html.Token

	// guards shared vars used for communication
	title	Title
	titles	chan Title
}

func Parse(r io.Reader) <-chan Title {
	l := &lexer{
		z: html.NewTokenizer(r),
		titles: make(chan Title),
		depth: 0,
	}
	go l.run()
	return l.titles
}

type lexerFunc func(*lexer) lexerFunc

func (l *lexer) run() {
	for f := lexNextToken; f != nil; {
		f = f(l)
	}
	close(l.titles)
}

func lexNextToken(l *lexer) lexerFunc {
	var tt = l.z.Next()
	l.t = l.z.Token()
	switch tt {
	case html.ErrorToken:
		if err := l.z.Err(); err != io.EOF {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		return nil
	case html.StartTagToken:
		l.depth++
		return lexTagToken
	case html.TextToken:
	case html.EndTagToken:
		l.depth--
		if l.depth == l.emitDepth {
			l.titles <- l.title
		}
	default:
	}
	return lexNextToken
}

// target is key value control of attrs mapped in lexer
type target struct {
	k, v string
	tt html.TokenType
	atom atom.Atom
	f lexerFunc
}

var targets []target

func lexTagToken(l *lexer) lexerFunc {
	switch l.t.DataAtom {
	case atom.Script, atom.Style:
		// prevent check
		return lexNextToken
	}
	return lexAttr
}

func lexNewTitle(l *lexer) lexerFunc {
	l.title = Title{}
	l.emitDepth = l.depth
	return lexNextToken
}

func lexAttr(l *lexer) lexerFunc {
	for _, a := range l.t.Attr {
		fmt.Printf("oi %s\n", a)
		for _, tv := range targets {
			if a.Key == tv.k && strings.Contains(a.Val, tv.v) {
				fmt.Println(tv)
				return tv.f
			}
		}
	}
	return lexNextToken
}

func lexTitleText(l *lexer) lexerFunc {
	for _, a := range l.t.Attr {
		if a.Key == "alt" {
			l.title.Text = a.Val
		}
	}
	return lexNextToken
}

func init() {
	targets = []target {
		target{
			k: "class",
			v: "lister-item mode-advanced",
			f: lexNewTitle,
		},
		target{
			k: "class",
			v: "loadlate",
			f: lexTitleText,
		},
	}
}