package temple

import (
	"fmt"
	"testing"
)

func TestCustomFunc(t *testing.T) {
	defer reset()
	AddFunc("greet", func(name string) string {
		return fmt.Sprintf("Hello, %s!", name)
	})
	// The test template calls the greet func
	if err := AddTemplate("test", `{{ greet "world"}}`); err != nil {
		t.Fatalf("Unexpected error in AddTemplate: %s", err.Error())
	}
	ExpectExecutorOutputs(t, Templates["test"], nil, "Hello, world!")
}
