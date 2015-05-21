// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"bytes"
	"testing"
)

// expectExecutorOutputs executes e and then adds an error to t if the output
// does not match expected.
func expectExecutorOutputs(t *testing.T, e Executor, data interface{}, expected string) {
	buf := bytes.NewBuffer([]byte{})
	if err := e.Execute(buf, data); err != nil {
		t.Errorf("Unexpected error executing template: %s", err.Error())
	}
	got := buf.String()
	if expected != got {
		t.Errorf("%T was not executed correctly. Expected `%s` but got `%s`.", e, expected, got)
	}
}
