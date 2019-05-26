package main

import (
	"bytes"
	"flag"
	"fmt"
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

	// os.Stdin is not a working ReadSeaker
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// When piping a newline is often appended, clean it
	stdin = bytes.TrimSpace(stdin)
	var frag = bytes.NewReader(stdin)

	err = ud.Replace(frag, *id, *file, *inplace, *replaceChild)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
