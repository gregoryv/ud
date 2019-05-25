package ud

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/workdir"
)

func TestReplace_errors(t *testing.T) {
	stdin := strings.NewReader("aaa")
	err := Replace(stdin, "a", "nosuchfile.html", true, true)
	if err == nil {
		t.Error("should fail when no file found")
	}

	// package global
	TempFile = func(string, string) (*os.File, error) {
		return nil, fmt.Errorf("oups")
	}
	wd, _ := workdir.TempDir()
	defer wd.RemoveAll()
	wd.WriteFile("index.html", []byte("<html></html>"))
	stdin = strings.NewReader("aaa")
	err = Replace(stdin, "a", wd.Join("index.html"), true, false)
	if err == nil {
		t.Error("should fail when temporary file cannot be created")
	}
	TempFile = ioutil.TempFile

}

func Test_findId(t *testing.T) {
	cases := []struct {
		in  string
		exp string
	}{
		{`<em id="who">github</em>`, "who"},
		{`<em>github</em>`, ""},
	}
	for _, c := range cases {
		r := strings.NewReader(c.in)
		got := findId(r)
		if got != c.exp {
			t.Error(got, c.exp)
		}
	}

}

func TestReplace(t *testing.T) {
	content := `<html><head><title></title></head>
<body><span id="a">
Hello, <em id="who">World</em>!
</span></body></html>`

	DefaultOutput = &discard{}

	file := "index.html"
	cases := []struct {
		id           string
		with         string
		exp          string
		inplace      bool
		replaceChild bool
	}{
		{
			id:           "a",
			with:         "aaa",
			exp:          "aaa",
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			with:         `<em id="who">github</em>`,
			exp:          `Hello, <em id="who">github</em>`,
			inplace:      true,
			replaceChild: false,
		},
		{
			id:           "",
			with:         `<em id="who">github</em>`,
			exp:          `Hello, <em id="who">github</em>`,
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			with:         `<em>github</em>`,
			exp:          `<em id="who">World</em>`,
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			with:         `<em id="who">github</em>`,
			exp:          `<em id="who">World</em>`,
			inplace:      false,
			replaceChild: false,
		},
	}
	for _, c := range cases {
		wd, _ := workdir.TempDir()
		wd.WriteFile(file, []byte(content))
		stdin := strings.NewReader(c.with)
		Replace(stdin, c.id, wd.Join(file), c.inplace, c.replaceChild)
		assert := asserter.New(t)
		got, _ := ioutil.ReadFile(wd.Join("index.html"))
		assert().Contains(got, c.exp)
		assert().Contains(got, "</body></html>")
		wd.RemoveAll()
	}
}

func TestNewInplaceWriter(t *testing.T) {
	w, err := NewInplaceWriter("x", ioutil.TempFile)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	w, err = NewInplaceWriter("x", func(string, string) (*os.File, error) {
		return nil, fmt.Errorf("oups")
	})
	if err == nil {
		t.Fatal(err)
	}
	if w != nil {
		w.Close()
	}
}

type discard struct{}

func (discard) Close() error { return nil }

func (discard) Write(b []byte) (int, error) { return len(b), nil }
