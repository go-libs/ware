// Package ware is a powerful package for easily create middleware layer in Golang.
//
// For a full guide visit https://github.com/futurespace/ware
//
// package main
//
// import (
//   "log"
//   "github.com/futurespace/ware"
// )
//
// func main() {
//   w := ware.New()
//
//   w.Use(func(c ware.Context, log *log.Logger) {
//     log.Println("before")
//     c.Next()
//     log.Println("after")
//   })
//
//   w.Run()
// }
package ware

import (
	"log"
	"os"
	"reflect"

	"github.com/codegangsta/inject"
)

// Middleware inject.Injector methods can be invoked to map services on a global level.
type Ware struct {
	inject.Injector
	handlers []Handler
	action   Handler
	logger   *log.Logger
}

// New creates a base bones Ware instance. Use this method if you want to have full control over the middleware that is used.
func New() *Ware {
	w := &Ware{Injector: inject.New(), action: func() {}, logger: log.New(os.Stdout, "[ware] ", 0)}
	w.Map(w.logger)
	w.Map(defaultReturnHandler())
	return w
}

// Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
// Will panic if any of the handlers is not a callable function
func (w *Ware) Handlers(handlers ...Handler) {
	w.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		w.Use(handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked. This is set to other Ware in a Ware.
func (w *Ware) Action(handler Handler) {
	validateHandler(handler)
	w.action = handler
}

// Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added
func (w *Ware) Use(handler Handler) {
	validateHandler(handler)

	w.handlers = append(w.handlers, handler)
}

// Run the ware.
func (w *Ware) Run() {
	w.CreateContext().Run()
}

// Creates a context.
func (w *Ware) CreateContext() *context {
	c := &context{inject.New(), w.handlers, w.action, nil, 0}
	c.SetParent(w)
	c.MapTo(c, (*Context)(nil))
	return c
}

// Handler can be any callable function. Martini attempts to inject services into the handler's argument list.
// Martini will panic if an argument could not be fullfilled via dependency injection.
type Handler interface{}

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("ware handler must be a callable func")
	}
}

// Context represents a request context. Services can be mapped on the request level from this interface.
type Context interface {
	inject.Injector
	// Next is an optional function that Middleware Handlers can call to yield the until after
	// the other Handlers have been executed. This works really well for any operations that must
	// happen after a request
	Next()
	// Written returns whether or not the response for this context has been written.
	Written() bool
	// The response instance.
	Out(interface{})
}

type context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	out      interface{}
	index    int
}

// Changes the response instance.
func (c *context) Out(o interface{}) {
	c.out = o
}

func (c *context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *context) Next() {
	c.index += 1
	c.Run()
}

func (c *context) Written() bool {
	out := reflect.ValueOf(c.out)
	if out.IsValid() {
		// Currently, defaults only mapping Written method.
		f := out.MethodByName("Written")
		if f.IsValid() && f.Kind() == reflect.Func {
			res := f.Call(nil)
			return res[0].Interface().(bool)
		}
	}
	return false
}

func (c *context) Run() {
	for c.index <= len(c.handlers) {
		_, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		// if the handler returned something, write it to the http response
		// TODO
		// if len(vals) > 0 {
		//   ev := c.Get(reflect.TypeOf(ReturnHandler(nil)))
		//   handleReturn := ev.Interface().(ReturnHandler)
		//   handleReturn(c, vals)
		// }

		if c.Written() {
			return
		}
	}
}
