package ud

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// HtmlRewriter replaces html tags or their content by id.
type HtmlRewriter struct {
	id       string
	child    bool
	fragment []byte
}

// NewHtmlRewriter returns a html rewriter. The rewriter can be used
// multiple times.
func NewHtmlRewriter(id string, child bool, fragment []byte) *HtmlRewriter {
	hr := &HtmlRewriter{
		id:       id,
		child:    child,
		fragment: fragment,
	}
	if hr.id == "" {
		hr.id = findId(fragment)
		hr.child = false // cannot be used when id is not given
	}
	return hr
}

func findId(frag []byte) string {
	z := html.NewTokenizer(bytes.NewReader(frag))
	for z.Next(); z.Err() != io.EOF; z.Next() {
		tok := z.Token()
		for _, attr := range tok.Attr {
			if attr.Key == "id" {
				return attr.Val
			}
		}
	}
	return ""
}

// Rewrite the incoming stream on r and write it to w.
func (hr *HtmlRewriter) Rewrite(w io.Writer, r io.Reader) error {
	z := html.NewTokenizer(r)
outer:
	for z.Next(); z.Err() != io.EOF; z.Next() {
		tok := z.Token()
		if idMatch(hr.id, tok) {
			if hr.child {
				fmt.Fprint(w, tok)
				skipChild(z)
				w.Write(hr.fragment)
				fmt.Fprint(w, "</", tok.Data, ">")
			} else {
				skip(z)
				w.Write(hr.fragment)
			}
			continue outer
		}
		fmt.Fprint(w, tok)
	}
	return nil
}

func idMatch(id string, tok html.Token) bool {
	for _, attr := range tok.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}
	return false
}

func skipChild(z *html.Tokenizer) html.Token {
	tt := z.Next()
	if tt != html.EndTagToken {
		return skip(z)
	}
	return z.Token()
}

func skip(z *html.Tokenizer) html.Token {
	depth := 1 // when 0 we stop
	for {
		tt := z.Next()
		switch tt {
		case html.StartTagToken:
			depth++
		case html.EndTagToken:
			depth--
		}
		if depth == 0 {
			return z.Token()
		}
	}
}
