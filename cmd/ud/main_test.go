package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/gregoryv/golden"
	"github.com/gregoryv/workdir"
)

var gobin = "go"

func init() {
	if runtime.GOOS == "windows" {
		gobin = "go.exe"
	}
	out, err := exec.Command(gobin, "build", ".").CombinedOutput()
	if err != nil {
		fmt.Println("Failed to build cmd/ud:", string(out))
		os.Exit(1)
	}
}

func TestCommand(t *testing.T) {
	wd, _ := workdir.TempDir()
	defer wd.RemoveAll()
	htmlFile, fragFile := setupFileAndFragment(wd)

	cmd := exec.Command("./ud", "-w", "-f", wd.Join(fragFile),
		"-html", wd.Join(htmlFile))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(out)
		t.Error(err)
	}

	got, _ := wd.Load(htmlFile)
	golden.Assert(t, string(got))
}

func TestMain(t *testing.T) {
	wd, _ := workdir.TempDir()
	defer wd.RemoveAll()
	htmlFile, fragFile := setupFileAndFragment(wd)

	Main("", wd.Join(htmlFile), wd.Join(fragFile), true, false,
		func(err error) {
			if err != nil {
				t.Error(err)
			}
		},
	)

	got, _ := wd.Load(htmlFile)
	golden.Assert(t, string(got))
}

func setupFileAndFragment(wd workdir.WorkDir) (htmlFile, fragFile string) {
	htmlFile = "index.html"
	content := []byte(`<html><body><h1 id="x">BIG</h1></body></html>`)
	wd.WriteFile(htmlFile, content)

	fragment := []byte(`<h2 id="x">small</h2>`)
	fragFile = "fragment.html"
	wd.WriteFile(fragFile, fragment)
	return
}
