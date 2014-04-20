package ware

import (
	"reflect"
)

// ReturnHandler is a service that Martini provides that is called
// when a route handler returns something. The ReturnHandler is
// responsible for writing to the ResponseWriter based on the values
// that are passed into this function.
type ReturnHandler func(Context, []reflect.Value)

func defaultReturnHandler() ReturnHandler {
	return func(ctx Context, vals []reflect.Value) {
		// TODO
	}
}
