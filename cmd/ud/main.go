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
	replaceChild := flag.Bool("c", false, "replace content not element")
	flag.Parse()

	// Fragments are usually small
	frag, err := ioutil.ReadAll(os.Stdin)
	fatal(err)

	// When piping a newline is often appended, clean it
	frag = bytes.TrimSpace(frag)

	var w io.WriteCloser = os.Stdout
	if *inplace {
		w, err = NewInplaceWriter(*file, TempFile)
		fatal(err)
	}
	defer w.Close()

	hr := ud.NewHtmlRewriter(*id, *replaceChild, frag)
	r, err := os.Open(*file)
	fatal(err)

	err = hr.Rewrite(w, r)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
