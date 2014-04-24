package ware

import (
	"log"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_Ware_Run(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run()
}

func Test_Ware_App(t *testing.T) {
	result := ""

	w := New()
	w.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	w.Use(func(c Context) {
		result += "bar"
		c.Next()
		result += "baz"
	})
	w.Action(func() {
		result += "bat"
	})

	w.Run()

	expect(t, result, "foobarbatbazban")
}

func Test_Ware_Handlers(t *testing.T) {
	result := ""

	batman := func(c Context) {
		result += "batman!"
	}

	w := New()
	w.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	w.Handlers(
		batman,
		batman,
		batman,
	)
	w.Action(func() {
		result += "bat"
	})

	w.Run()

	expect(t, result, "batman!batman!batman!bat")
}

func Test_Ware_Logger_SetPrefix(t *testing.T) {
	prefix := ""

	w := New()
	w.Use(func(log *log.Logger) {
		log.SetPrefix("[martini]")
	})
	w.Use(func(log *log.Logger) {
		prefix = log.Prefix()
	})
	w.Run()

	expect(t, prefix, "[martini]")
}
