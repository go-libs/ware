# Ware [![wercker status](https://app.wercker.com/status/1569ebfba816e02d463a2b55b2000744/s/ "wercker status")](https://app.wercker.com/project/bykey/1569ebfba816e02d463a2b55b2000744) [![GoDoc](https://godoc.org/github.com/futurespace/ware?status.png)](https://godoc.org/github.com/futurespace/ware)

Easily create middleware layer in Golang.   
**Forked from [martini][]**.   
Dependence [inject][] package.


## Getting Started

```go
// ware.go
package main

import (
        "log"
        "github.com/futurespace/ware"
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
go get github.com/futurespace/ware
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
        . "github.com/futurespace/ware"
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

Mapping values to interface

```go
type Deploy interface {
        Do()
}
type deploy struct{}
func (d *deploy) Do() {}
func NewDeploy() Deploy {
        return &deploy{}
}
type Compress struct {
        inject.Injector
        *Ware
        d Deploy
}
func NewCompress() *Compress {
        w := New()
        d := NewDeploy()
        w.MapTo(d, (*Deploy)(nil))
        c := &Compress{inject.New(), w, d}
        c.Map(w)
        return c
}


c := NewCompress()
c.Use(func(d Deploy) {
        fmt.Println(d)
})
c.Run()
```

Sets the output prefix for the logger. (Default `[ware]`)

```go
w.Use(func(log *log.Logger) {
        log.SetPrefix("[martini]")
})
```


## API

#### New() 
> Creates a Ware instance.

#### Ware.Handlers(handlers ...Handler)
> Sets middlewares.

#### Ware.Action(Handler)
> Sets the handler that will be called after all the middleware has been invoked.

#### Ware.Use(Handler)
> Adds a middleware.

#### Ware.Run()
> Invokes the ware app.

#### Ware.CreateContext()
> Creates a new context.   
> ***NOTE***: In martini, this api is private, but sometime we need to hack the context!

#### Context.Out(o interface{})
> Sets response instance.

#### Context.Next()

#### Context.Written()
> Stops to invoke next middleware handler, the response instance responded.

#### Context.Run()
> Invokes the context.

### Others

  See [inject][].


## License

MIT


[martini]: https://github.com/go-martini/martini
[inject]: github.com/codegangsta/inject
