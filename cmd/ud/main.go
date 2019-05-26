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

	var frag ud.Fragment = os.Stdin
	if *id == "" {
		// If no Id is given the frag must be a working ReadSeeker
		// os.Stdin is Not
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		frag = bytes.NewReader(stdin)
	}
	err := ud.Replace(frag, *id, *file, *inplace, *replaceChild)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
