// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package temple

import (
	"fmt"
	"testing"
)

func TestCustomFunc(t *testing.T) {
	g := NewGroup()
	g.AddFunc("greet", func(name string) string {
		return fmt.Sprintf("Hello, %s!", name)
	})
	// The test template calls the greet func
	if err := g.AddTemplate("test", `{{ greet "world"}}`); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	expectExecutorOutputs(t, g.templates["test"], nil, "Hello, world!")
}
