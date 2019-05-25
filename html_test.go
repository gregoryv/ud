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
	wd, _ := workdir.TempDir()
	content := `<html><head><title></title></head>
<body><span id="a">
Hello, <em id="who">World</em>!
</span></body></html>`
	defer wd.RemoveAll()
	wd.WriteFile("index.html", []byte(content))
	stdin := strings.NewReader("aaa")
	err := Replace(stdin, "a", wd.Join("index.html"), true, true)
	assert := asserter.New(t)
	assert(err == nil).Fatal(err)
	got, _ := ioutil.ReadFile(wd.Join("index.html"))
	assert().Contains(got, "aaa")
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
