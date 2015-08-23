package goldi

import (
	"fmt"
	"reflect"
)

type funcType struct {
	function interface{}
}

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
