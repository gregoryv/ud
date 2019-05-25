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

func TestReplace(t *testing.T) {
	content := `<html><head><title></title></head>
<body><span id="a">
Hello, <em id="who">World</em>!
</span></body></html>`

	file := "index.html"
	cases := []struct {
		id           string
		with         string
		inplace      bool
		replaceChild bool
	}{
		{"a", "aaa", true, true},
	}
	for _, c := range cases {
		wd, _ := workdir.TempDir()
		wd.WriteFile(file, []byte(content))
		stdin := strings.NewReader(c.with)
		err := Replace(stdin, c.id, wd.Join(file), c.inplace, c.replaceChild)
		assert := asserter.New(t)
		assert(err == nil).Error(err)
		got, _ := ioutil.ReadFile(wd.Join("index.html"))
		assert().Contains(got, c.with)
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
