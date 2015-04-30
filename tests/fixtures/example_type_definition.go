package fixtures

import (
	. "github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var types = NewTypeRegistry()

func init() {
	types.RegisterType("goldi.test.foo", testAPI.NewFoo)
	types.RegisterType("goldi.test.bar", testAPI.NewBar)

	types.Register("goldi.test.baz", NewType(testAPI.NewBaz, "Hello", "World!"))
}
