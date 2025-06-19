package di_test

import (
	"testing"

	"github.com/kimnguyenlong/go-di"
)

type A struct {
	Name string
}

type B struct {
	Name string
	A    *A `di:"a"`
}

func TestInject(t *testing.T) {
	a := &A{Name: "A"}
	b := &B{Name: "B"}

	// Register the objects
	di.Plug(
		&di.Object{Name: "a", Value: a},
		&di.Object{Name: "b", Value: b},
	)

	// Wire the dependencies
	if err := di.Wire(); err != nil {
		t.Fatalf("Failed to wire dependencies: %v", err)
	}

	if b.A == nil {
		t.Fatal("Expected B.A to be injected, but it is nil")
	}

	if b.A.Name != "A" {
		t.Fatalf("Expected B.A.Name to be 'A', got '%s'", b.A.Name)
	}

	if b.A != a {
		t.Fatalf("Expected B.A to be the same instance as A, but they are different")
	}
}
