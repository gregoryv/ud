package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Used to create temporary files for writing inplace
var (
	TempFile TempFiler = func(dir, pattern string) (*os.File, error) {
		return os.Create(filepath.Join(dir, pattern+".tmp"))
	}
	DefaultOutput io.WriteCloser = os.Stdout
)

type TempFiler func(string, string) (*os.File, error)

func NewInplaceWriter(file string, newTemp TempFiler) (*InplaceWriter, error) {
	tmp, err := newTemp("", "ud")
	if err != nil {
		return nil, err
	}
	return &InplaceWriter{tmp, file}, nil
}

type InplaceWriter struct {
	tmp  *os.File
	dest string
}

func (w *InplaceWriter) Write(b []byte) (int, error) {
	return w.tmp.Write(b)
}

func (w *InplaceWriter) Close() error {
	w.tmp.Close()
	os.Chmod(w.tmp.Name(), 0644)
	log.Println("rewrite", w.tmp.Name(), "to", w.dest)
	return os.Rename(w.tmp.Name(), w.dest)
}
