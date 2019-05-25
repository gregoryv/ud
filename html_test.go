package ud

import (
	"bytes"
	"io/ioutil"
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

	file := "index.html"
	cases := []struct {
		id           string
		with         string
		exp          string
		inplace      bool
		replaceChild bool
	}{
		{"a", "aaa", "aaa", true, true},
		{"", `<em id="who">github</em>`, `Hello, <em id="who">github</em>`, true, false},
		{"", `<em id="who">github</em>`, `Hello, <em id="who">github</em>`, true, true},
		{"", `<em>github</em>`, `<em id="who">World</em>`, true, true},
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

func TestTraverse(t *testing.T) {
	content := `<html><head><title></title></head>
<body><span id="a">hepp</span></body></html>`

	doc := strings.NewReader(content)
	buf := bytes.NewBufferString("")
	replace(doc, "a", strings.NewReader("hello"), buf, false)
	assert := asserter.New(t)
	got := buf.String()
	assert().Contains(got, "hello")
	if strings.Index(got, "span") != -1 {
		t.Error("found span")
	}

	buf = bytes.NewBufferString("")
	doc = strings.NewReader(content)
	replace(doc, "a", strings.NewReader("hello"), buf, true)
	got = buf.String()
	assert().Contains(got, "hello</span>")

}
