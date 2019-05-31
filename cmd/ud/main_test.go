package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
