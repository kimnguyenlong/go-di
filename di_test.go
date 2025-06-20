package di_test

import (
	"testing"

	"github.com/kimnguyenlong/go-di"
)

type A interface {
	DoA() string
}

type B interface {
	DoB() string
}

type AImpl struct{}

func (a *AImpl) DoA() string {
	return "A implementation"
}

type BImpl struct {
	A A
}

func (b *BImpl) DoB() string {
	return "B implementation"
}

func TestSimpleInject(t *testing.T) {
	c := di.NewContainer()
	app := &struct {
		A A `di:"a"`
		B B `di:"b"`
	}{}
	c.Plug("app", app)
	c.Plug("a", &AImpl{})
	c.Plug("b", &BImpl{})
	if err := c.Wire(); err != nil {
		t.Fatalf("di.Wire() failed: %v", err)
	}
	if app.A == nil {
		t.Fatal("app.A is nil after injection")
	}
	if app.B == nil {
		t.Fatal("app.B is nil after injection")
	}
	if app.A.DoA() != "A implementation" {
		t.Errorf("app.A.DoA() = %s, want 'A implementation'", app.A.DoA())
	}
	if app.B.DoB() != "B implementation" {
		t.Errorf("app.B.DoB() = %s, want 'B implementation'", app.B.DoB())
	}
}

func TestNestedInject(t *testing.T) {
	c := di.NewContainer()
	type Nested struct {
		A A `di:"a"`
	}
	type App struct {
		Nested *Nested
	}
	app := &App{Nested: &Nested{}}
	c.Plug("app", app)
	c.Plug("a", &AImpl{})
	if err := c.Wire(); err != nil {
		t.Fatalf("di.Wire() failed: %v", err)
	}
	if app.Nested.A == nil {
		t.Fatal("app.Nested.A is nil after injection")
	}
	if app.Nested.A.DoA() != "A implementation" {
		t.Errorf("app.Nested.A.DoA() = %s, want 'A implementation'", app.Nested.A.DoA())
	}
}
func TestMissingDependency(t *testing.T) {
	c := di.NewContainer()
	app := &struct {
		A A `di:"a"`
	}{}
	c.Plug("app", app)
	// Do not plug "a"
	err := c.Wire()
	if err == nil {
		t.Fatal("expected error when dependency is missing, got nil")
	}
}

func TestTypeMismatch(t *testing.T) {
	c := di.NewContainer()
	app := &struct {
		A A `di:"a"`
	}{}
	c.Plug("app", app)
	c.Plug("a", "not an AImpl")
	err := c.Wire()
	if err == nil {
		t.Fatal("expected error on type mismatch, got nil")
	}
}

func TestNoDIField(t *testing.T) {
	c := di.NewContainer()
	app := &struct {
		X int
	}{X: 42}
	c.Plug("app", app)
	if err := c.Wire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.X != 42 {
		t.Errorf("expected X to remain 42, got %d", app.X)
	}
}

func TestNilContainer(t *testing.T) {
	var c *di.Container = nil
	err := c.Wire()
	if err == nil {
		t.Fatal("expected error when container is nil, got nil")
	}
}

func TestNilObject(t *testing.T) {
	c := di.NewContainer()
	c.Plug("nilobj", nil)
	err := c.Wire()
	if err == nil {
		t.Fatal("expected error when object is nil, got nil")
	}
}

func TestNonPointerStruct(t *testing.T) {
	c := di.NewContainer()
	obj := struct {
		A A `di:"a"`
	}{}
	c.Plug("obj", obj) // Not a pointer
	c.Plug("a", &AImpl{})
	err := c.Wire()
	if err == nil {
		t.Fatal("expected error for non-pointer struct, got nil")
	}
}
