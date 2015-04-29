package fixtures

import (
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var typeDef = goldi.NewAppDefinition()

func init() {
	typeDef.RegisterType("goldi.test.foo", testAPI.NewFoo)
	typeDef.RegisterType("goldi.test.bar", testAPI.NewBar)
	typeDef.RegisterType("goldi.test.baz", testAPI.NewBaz)
}
