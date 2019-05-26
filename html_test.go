package ud

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

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
	DefaultOutput = &discard{}

	file := "index.html"
	cases := []struct {
		doc          string
		id           string
		frag         string
		exp          string
		replaceChild bool
	}{
		{
			doc:          `<b><i id="x"></i></b>`, // empty start
			id:           "x",
			frag:         `content`, // ok id
			exp:          `<b><i id="x">content</i></b>`,
			replaceChild: true,
		},
		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "",
			frag:         `<em id="x">A</em>`, // ok id
			exp:          `<b><em id="x">A</em></b>`,
			replaceChild: false,
		},
		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "",
			frag:         `<i id="X"></i>`, // wrong id
			exp:          `<b><i id="x">a</i></b>`,
			replaceChild: false,
		},
		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "x",
			frag:         "",
			exp:          `<b></b>`,
			replaceChild: false,
		},
		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "", // no id given
			frag:         "", // no id found
			exp:          `<b><i id="x">a</i></b>`,
			replaceChild: false,
		},
		{
			doc:          `<b><i id="x"><span>here</span></i></b>`,
			id:           "x",
			frag:         "",
			exp:          `<b></b>`,
			replaceChild: false,
		},

		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "x",
			frag:         `A`,
			exp:          `<b><i id="x">A</i></b>`,
			replaceChild: true,
		},
		{
			doc:          `<b><i id="x">a</i></b>`,
			id:           "", // empty
			frag:         `<i id="x">A</i>`,
			exp:          `<b><i id="x">A</i></b>`,
			replaceChild: true, // no effect when no id is given
		},
	}
	for _, c := range cases {
		wd, _ := workdir.TempDir()
		wd.WriteFile(file, []byte(c.doc))
		frag := strings.NewReader(c.frag)
		Replace(frag, c.id, wd.Join(file), true, c.replaceChild)
		got, _ := ioutil.ReadFile(wd.Join("index.html"))
		if string(got) != c.exp {
			t.Log("doc.:", c.doc)
			t.Logf("id..: %q, child: %v", c.id, c.replaceChild)
			t.Log("frag:", c.frag)
			t.Log("got.:", string(got))
			t.Log("exp.:", c.exp)

			t.Log()
			t.Fail()
		}
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

func Test_getOutput(t *testing.T) {
	got, _ := getOutput(false, "")
	if got != DefaultOutput {
		t.Fail()
	}
}

type discard struct{}

func (discard) Close() error { return nil }

func (discard) Write(b []byte) (int, error) { return len(b), nil }

// badTemper fails to create a temporary file, love the name :-)
func badTemper(string, string) (*os.File, error) {
	return nil, fmt.Errorf("oups")
}
