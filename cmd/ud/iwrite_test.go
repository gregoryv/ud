package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewInplaceWriter(t *testing.T) {
	cases := []struct {
		fn  TempFiler
		exp bool
	}{
		{
			fn:  ioutil.TempFile,
			exp: true,
		},
		{
			fn:  badTemper,
			exp: false,
		},
	}
	for _, c := range cases {
		fh, _ := ioutil.TempFile("", "x")
		fh.Close()

		w, err := NewInplaceWriter(fh.Name(), c.fn)
		if c.exp != (err == nil) {
			t.Error(err)
		}
		if w != nil {
			w.Close()
		}
		os.RemoveAll(fh.Name())
	}
}

// badTemper fails to create a temporary file, love the name :-)
func badTemper(string, string) (*os.File, error) {
	return nil, fmt.Errorf("oups")
}
