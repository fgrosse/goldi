package goldi
import (
	"fmt"
	"reflect"
)

type FuncType struct {
	function interface{}
}

func NewFuncType(function interface{}) *FuncType {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not register func type: %v", r))
		}
	}()

	structType := reflect.TypeOf(function)
	if structType.Kind() != reflect.Func {
		panic(fmt.Errorf("the given type must be a function (given %T)", function))
	}

	return &FuncType{function}
}

func (t *FuncType) Arguments() []interface{} {
	return []interface{}{}
}

func (t *FuncType) Generate(parameterResolver *ParameterResolver) interface{} {
	if t.function == nil {
		panic(fmt.Errorf("could not generate type: this func type is not initialized. Did you use NewFuncType to create it?"))
	}

	return t.function
}
