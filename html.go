package ud

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

func Replace(frag Fragment, id, file string, inplace, child bool) (err error) {
	if id == "" {
		id = findId(frag)
		child = false // cannot be used when id is not given
	}
	if id == "" {
		return fmt.Errorf("No id specified")
	}
	doc, err := os.Open(file)
	if err != nil {
		return
	}
	defer doc.Close()
	out, err := getOutput(inplace, file)
	if err != nil {
		return
	}
	replace(doc, frag, id, out, child)
	return out.Close()
}

type Fragment interface {
	io.ReadSeeker
}

// Used to create temporary files for writing inplace
var TempFile TempFiler = ioutil.TempFile
var DefaultOutput io.WriteCloser = os.Stdout

// getOutput returns writer, caller must call Close when done.
func getOutput(inplace bool, file string) (io.WriteCloser, error) {
	if inplace {
		return NewInplaceWriter(file, TempFile)
	}
	return DefaultOutput, nil
}

type InplaceWriter struct {
	tmp  *os.File
	dest string
}

type TempFiler func(string, string) (*os.File, error)

func NewInplaceWriter(file string, newTemp TempFiler) (*InplaceWriter, error) {
	tmp, err := newTemp("", "ud")
	if err != nil {
		return nil, err
	}
	return &InplaceWriter{tmp, file}, nil
}

func (w *InplaceWriter) Write(b []byte) (int, error) {
	return w.tmp.Write(b)
}

func (w *InplaceWriter) Close() error {
	w.tmp.Close()
	return os.Rename(w.tmp.Name(), w.dest)
}

func findId(frag Fragment) string {
	defer frag.Seek(0, 0) // reset to beginning
	z := html.NewTokenizer(frag)
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

func replace(doc, frag io.Reader, id string, w io.Writer, child bool) {
	z := html.NewTokenizer(doc)
outer:
	for z.Next(); z.Err() != io.EOF; z.Next() {
		tok := z.Token()
		for _, attr := range tok.Attr {
			if idMatch(id, attr) {
				if child {
					fmt.Fprint(w, tok)
					tt := z.Next()
					if tt != html.EndTagToken {
						skip(z)
					}
					io.Copy(w, frag)
					fmt.Fprint(w, "</", tok.Data, ">")
				} else {
					skip(z)
					io.Copy(w, frag)
				}

				continue outer
			}
		}
		fmt.Fprint(w, tok)
	}
}

func idMatch(id string, attr html.Attribute) bool {
	return attr.Key == "id" && attr.Val == id
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
