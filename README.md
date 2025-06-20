## go-di

Simple dependency injection framework for Go applications.

### Feature Overview

- Register any object instance in a central container with a unique name.
- Inject dependencies into struct fields using the `di` struct tag.
- Supports recursive injection for nested struct pointers.
- Type safety checks during injection to prevent mismatched assignments.

## Guide

### Installation

```sh
// go get github.com/kimnguyenlong/go-di@{version}
go get github.com/kimnguyenlong/go-di@latest
```

### Example

```go
package main

import (
	"log"

	"github.com/kimnguyenlong/go-di"
)

func main() {
	container := di.NewContainer()             // init empty container
	app := &App{Dependencies: &Dependencies{}} // Object A hasn't been injected yet
	objA := &AImpl{}                           // init object A

	container.Plug("app", app)               // register app as name "app"
	container.Plug("aa", objA)               // register objA as name "a"
	if err := container.Wire(); err != nil { // err is nil
		log.Fatalf("write failed: %s", err)
	}

	app.Dependencies.A.DoA() // print "A implementation"
}

type App struct {
	Dependencies *Dependencies
}

type Dependencies struct {
	A A `di:"a"` // will be injected with object named "a"
}

type A interface {
	DoA()
}

type AImpl struct{}

func (a *AImpl) DoA() {
	log.Println("A implementation")
}
```

## License

[MIT](https://github.com/kimnguyenlong/go-di/blob/main/LICENSE)