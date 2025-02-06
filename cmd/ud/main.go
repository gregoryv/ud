package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gregoryv/ud"
)

//go:generate stamp -clfile ../../CHANGELOG.md -go build_stamp.go
func main() {
	id := flag.String("i", "", "Id of element")
	file := flag.String("html", "", "html file to modify")
	inplace := flag.Bool("w", false, "writes to file inplace")
	fragFile := flag.String("f", "", "fragment to use")
	child := flag.Bool("c", false, "replace content not element")
	v := flag.Bool("v", false, "Print version and exit")
	vv := flag.Bool("vv", false, "Print version with details and exit")
	flag.Parse()

	if *v {
		fmt.Println(ud.Version())
		os.Exit(0)
	}
	if *vv {
		fmt.Printf("%s-%s\n", ud.Version(), ud.Revision())
		os.Exit(0)
	}

	err := Main(*id, *file, *fragFile, *inplace, *child)
	logError(err)
	os.Exit(exitCode(err))
}

func Main(id, file, fragFile string, inplace, child bool) error {

	frag, err := readFragment(fragFile)
	if err != nil {
		return err
	}

	// When piping a newline is often appended, clean it
	frag = bytes.TrimSpace(frag)
	w, err := newWriteCloser(inplace, file)
	if err != nil {
		return err
	}

	hr := ud.NewHtmlRewriter(id, child, frag)
	r, err := os.Open(file)
	if err != nil {
		return err
	}

	err = hr.Rewrite(w, r)
	if err != nil {
		return err
	}
	return w.Close()
}

func newWriteCloser(inplace bool, file string) (io.WriteCloser, error) {
	if inplace {
		return NewInplaceWriter(file, TempFile)
	}
	return os.Stdout, nil
}

func readFragment(filename string) ([]byte, error) {
	if filename == "" {
		return ioutil.ReadAll(os.Stdin)
	}
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return ioutil.ReadAll(fh)
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func exitCode(err error) ExitCode {
	if err != nil {
		return ExitFail
	}
	return ExitOk
}

type ExitCode = int

const (
	ExitOk ExitCode = iota
	ExitFail
)
