// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"os/exec"
	"testing"
)

const (
	destFile = "test_files/templates.go"
	runFile  = "test_files/run.go"
)

func TestBuild(t *testing.T) {
	// Generate a go source file with build
	if err := Build("test_files/templates", destFile, "test_files/partials", "test_files/layouts", "main"); err != nil {
		t.Error(err)
	}
	// Use go run to run the file together with the run file
	cmd := exec.Command("go", "run", destFile, runFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	expected := "<html><head><title>Todos</title></head><body><ul><li>One</li><li>Two</li><li>Three</li></ul></body></html>"
	if string(output) != expected {
		t.Errorf("Output from generated code was not correct.\nExpected %s\nBut got:  %s", expected, string(output))
	}
}
