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

func main() {
	id := flag.String("i", "", "Id of element")
	file := flag.String("html", "", "html file to modify")
	inplace := flag.Bool("w", false, "writes to file inplace")
	fragFile := flag.String("f", "", "fragment to use")
	child := flag.Bool("c", false, "replace content not element")
	flag.Parse()

	Main(*id, *file, *fragFile, *inplace, *child, fatal)
}

func Main(id, file, fragFile string, inplace, child bool, handle func(error)) {
	var (
		err    error
		frag   []byte    // Fragments are usually small
		fragIn io.Reader = os.Stdin
	)

	if fragFile != "" {
		fragIn, err = os.Open(fragFile)
		handle(err)
	}
	frag, err = ioutil.ReadAll(fragIn)
	handle(err)

	// When piping a newline is often appended, clean it
	frag = bytes.TrimSpace(frag)

	var w io.WriteCloser = os.Stdout
	if inplace {
		w, err = NewInplaceWriter(file, TempFile)
		handle(err)
	}
	defer w.Close()

	hr := ud.NewHtmlRewriter(id, child, frag)
	r, err := os.Open(file)
	handle(err)

	err = hr.Rewrite(w, r)
	handle(err)
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
