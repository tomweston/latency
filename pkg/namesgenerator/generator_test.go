package namesgenerator

import (
	"testing"
	"time"
)

func TestNameGenerator(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	_ = NewNameGenerator(seed)
}

func TestNameGenerator_Generate(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := NewNameGenerator(seed)

	name := nameGenerator.Generate()
	if name == "" {
		t.Fatalf("Expected a new name but got a blank string")
	}
}
