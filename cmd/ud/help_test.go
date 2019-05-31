package main

import (
	"os/exec"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/workdir"
)

func TestHelp(t *testing.T) {
	out, err := exec.Command("./ud", "-h").CombinedOutput()
	if err.Error() != "exit status 2" {
		t.Error(string(out), err)
	}
	assert := asserter.New(t)
	assert().Contains(out, "Usage of")
}

func TestBasicOperation(t *testing.T) {
	wd, _ := workdir.TempDir()

	htmlFile := "index.html"
	content := []byte(`<html><body><h1 id="x">BIG</h1></body></html>`)
	wd.WriteFile(htmlFile, content)

	fragment := []byte(`<h2 id="x">small</h2>`)
	fragFile := "fragment.html"
	wd.WriteFile(fragFile, fragment)

	cmd := exec.Command("./ud", "-w", "-f", wd.Join(fragFile),
		"-html", wd.Join(htmlFile))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(out))
	}
	newContent, err := wd.Load(htmlFile)
	if err != nil {
		t.Error(err)
	}
	got := string(newContent)
	exp := `<html><body><h2 id="x">small</h2></body></html>`
	if got != exp {
		t.Error(got)
	}
	wd.RemoveAll()
}
