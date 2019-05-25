package ud

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/net/html"
)

func Replace(r io.ReadSeeker, id, file string, inplace, replaceChild bool) (err error) {
	fh, err := os.Open(file)
	if err != nil {
		return
	}
	if id == "" {
		id = findId(r)
		replaceChild = false // cannot be used when id is not given
	}
	if id == "" {
		return fmt.Errorf("No id specified")
	}
	if !inplace {
		replace(fh, r, id, os.Stdout, replaceChild)
		return
	}
	out, err := ioutil.TempFile("", path.Base(file))
	if err != nil {
		return
	}
	replace(fh, r, id, out, replaceChild)
	out.Close()
	return os.Rename(out.Name(), file)
}

func findId(r io.ReadSeeker) string {
	defer r.Seek(0, 0)
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		tok := z.Token()
		switch tt {
		case html.ErrorToken:
			break
		case html.StartTagToken:
			for _, attr := range tok.Attr {
				if attr.Key == "id" {
					return attr.Val
				}
			}
		}

		if z.Err() == io.EOF {
			break
		}
	}
	return ""
}

func replace(doc, newContent io.Reader, id string, w io.Writer, c bool) {
	z := html.NewTokenizer(doc)
	emitToken := func(t html.Token) {
		fmt.Fprint(w, t)
	}
outer:
	for {
		tt := z.Next()
		tok := z.Token()
		switch tt {
		case html.ErrorToken:
			break
		case html.StartTagToken:
			for _, attr := range tok.Attr {
				if attr.Key == "id" && attr.Val == id {
					if c {
						emitToken(tok)
						z.Next()
						skip(z)
						io.Copy(w, newContent)
						fmt.Fprint(w, "</", tok.Data, ">")
					} else {
						skip(z)
						io.Copy(w, newContent)
					}

					continue outer
				}
			}
		}
		emitToken(tok)

		if z.Err() == io.EOF {
			break
		}
	}
}

func skip(z *html.Tokenizer) {
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
			break
		}
	}
}
