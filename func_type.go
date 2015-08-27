package goldi

import (
	"fmt"
	"reflect"
)

type funcType struct {
	function interface{}
}

// NewFuncType creates a new TypeFactory that will return a method value
//
// Goldigen yaml syntax example:
//     my_func_type:
//         package: github.com/fgrosse/foobar
//         func:    DoStuff
func NewFuncType(function interface{}) TypeFactory {
	structType := reflect.TypeOf(function)
	if structType.Kind() != reflect.Func {
		return newInvalidType(fmt.Errorf("the given type must be a function (given %T)", function))
	}

	return &funcType{function}
}

func (t *funcType) Arguments() []interface{} {
	return []interface{}{}
}

func (t *funcType) Generate(parameterResolver *ParameterResolver) (interface{}, error) {
	return t.function, nil
}
