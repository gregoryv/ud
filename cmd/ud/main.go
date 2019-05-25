package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gregoryv/ud"
)

func main() {
	id := flag.String("i", "", "Id of element")
	file := flag.String("html", "", "html file to modify")
	inplace := flag.Bool("w", false, "writes to file inplace")
	replaceChild := flag.Bool("c", false, "replace content not element")
	flag.Parse()

	err := ud.Replace(os.Stdin, *id, *file, *inplace, *replaceChild)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
