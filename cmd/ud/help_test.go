package main

import (
	"os/exec"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestHelp(t *testing.T) {
	out, err := exec.Command("./ud", "-h").CombinedOutput()
	if err != nil && err.Error() != "exit status 2" {
		t.Error(string(out), err)
	}
	assert := asserter.New(t)
	assert().Contains(out, "Usage of")
}
