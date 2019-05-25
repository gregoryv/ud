package ud

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

func Replace(r io.ReadSeeker, id, file string,
	inplace, replaceChild bool) (err error) {
	if id == "" {
		id = findId(r)
		replaceChild = false // cannot be used when id is not given
	}
	if id == "" {
		return fmt.Errorf("No id specified")
	}
	fh, err := os.Open(file)
	if err != nil {
		return
	}
	out, err := getOutput(inplace, file)
	if err != nil {
		return
	}
	replace(fh, r, id, out, replaceChild)
	return out.Close()
}

// Used to create temporary files for writing inplace
var TempFile = ioutil.TempFile

// getOutput returns writer, caller must call Close when done.
func getOutput(inplace bool, file string) (io.WriteCloser, error) {
	if inplace {
		return NewInplaceWriter(file, TempFile)
	}
	return os.Stdout, nil
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

func findId(r io.ReadSeeker) string {
	defer r.Seek(0, 0)
	z := html.NewTokenizer(r)
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

func replace(doc, r io.Reader, id string, w io.Writer,
	replaceChild bool) {
	z := html.NewTokenizer(doc)
outer:
	for z.Next(); z.Err() != io.EOF; z.Next() {
		tok := z.Token()
		for _, attr := range tok.Attr {
			if idMatch(id, attr) {
				if replaceChild {
					fmt.Fprint(w, tok)
					z.Next()
					skip(z)
					io.Copy(w, r)
					fmt.Fprint(w, "</", tok.Data, ">")
				} else {
					skip(z)
					io.Copy(w, r)
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
