# Ware [![wercker status](https://app.wercker.com/status/1569ebfba816e02d463a2b55b2000744/s/ "wercker status")](https://app.wercker.com/project/bykey/1569ebfba816e02d463a2b55b2000744) [![GoDoc](https://godoc.org/github.com/futurespaceio/ware?status.png)](https://godoc.org/github.com/futurespaceio/ware)

Easily create middleware layer in Golang.   
Forked from [martini][].   
Dependence the [inject][] package.


## Getting Started

```go
// ware.go
package main

import (
  "log"
  "github.com/futurespaceio/ware"
)

func main() {
  w := ware.New()

  w.Use(func(c ware.Context, log *log.Logger) {
    log.Println("before")
    c.Next()
    log.Println("after")
  })

  w.Run()
}
```

Install the Ware package:

```
go get github.com/futurespaceio/ware
```

Run test:

```
go run ware.go
```

Compose Wares:

```go
package main

import (
  "log"

  "github.com/codegangsta/inject"
  . "github.com/futurespaceio/ware"
)

type Builder struct {
  inject.Injector
  *Ware
}

func NewBuilder() *Builder {
  w := New()
  b := &Builder{inject.New(), w}
  b.Map(w)
  return b
}

type Packer struct {
  inject.Injector
  *Ware
}

func (p *Packer) Handle() {
  p.Run()
}

func NewPacker() *Packer {
  w := New()
  p := &Packer{inject.New(), w}
  p.Map(w)
  return p
}
```

```go
b := NewBuilder()
b.Use(func (log *log.Logger) {
  log.Println("build...")
})
b.Run()

// Compose other Ware
p := NewPacker()
p.Use(func (log *log.Logger) {
  log.Println("pack...")
})
b.Action(p.Handle)
b.Run()
```


## API

### New() 

### Ware.Handlers(handlers ...Handler)

### Ware.Action(Handler)

### Ware.Use(Handler)

### Ware.Run(Handler)

### Others

  See [inject][].


## License

MIT


[martini]: https://github.com/go-martini/martini
[inject]: github.com/codegangsta/inject
