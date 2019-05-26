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
	TempFile = badTemper

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

<map id="g_all"></map>
</span></body></html>`

	DefaultOutput = &discard{}

	file := "index.html"
	cases := []struct {
		id           string
		frag         string
		exp          string
		inplace      bool
		replaceChild bool
	}{
		{
			id:           "",
			frag:         `<map id="g_all">something</map>`,
			exp:          `<map id="g_all">something</map>`,
			inplace:      true,
			replaceChild: false,
		},
		{
			id:           "a",
			frag:         "aaa",
			exp:          "aaa",
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			frag:         `<em id="who">github</em>`,
			exp:          `Hello, <em id="who">github</em>`,
			inplace:      true,
			replaceChild: false,
		},
		{
			id:           "",
			frag:         `<em id="who">github</em>`,
			exp:          `Hello, <em id="who">github</em>`,
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			frag:         `<em>github</em>`,
			exp:          `<em id="who">World</em>`,
			inplace:      true,
			replaceChild: true,
		},
		{
			id:           "",
			frag:         `<em id="who">github</em>`,
			exp:          `<em id="who">World</em>`,
			inplace:      false,
			replaceChild: false,
		},
	}
	for _, c := range cases {
		wd, _ := workdir.TempDir()
		wd.WriteFile(file, []byte(content))
		frag := strings.NewReader(c.frag)
		Replace(frag, c.id, wd.Join(file), c.inplace, c.replaceChild)
		assert := asserter.New(t)
		got, _ := ioutil.ReadFile(wd.Join("index.html"))
		assert().Contains(got, c.exp)
		assert().Contains(got, "</body></html>")
		wd.RemoveAll()
	}
}

func TestNewInplaceWriter(t *testing.T) {
	cases := []struct {
		fn  TempFiler
		exp bool
	}{
		{
			fn:  ioutil.TempFile,
			exp: true,
		},
		{
			fn:  badTemper,
			exp: false,
		},
	}
	for _, c := range cases {
		fh, _ := ioutil.TempFile("", "x")
		fh.Close()

		w, err := NewInplaceWriter(fh.Name(), c.fn)
		if c.exp != (err == nil) {
			t.Error(err)
		}
		if w != nil {
			w.Close()
		}
		os.RemoveAll(fh.Name())
	}
}

type discard struct{}

func (discard) Close() error { return nil }

func (discard) Write(b []byte) (int, error) { return len(b), nil }

// badTemper fails to create a temporary file, love the name :-)
func badTemper(string, string) (*os.File, error) {
	return nil, fmt.Errorf("oups")
}
