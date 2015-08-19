package goldi

import (
	"fmt"
	"reflect"
	"unicode"
)

type FuncReferenceType struct {
	TypeID       string
	FunctionName string
}

func NewFuncReferenceType(typeID, functionName string) *FuncReferenceType {
	if functionName == "" || unicode.IsLower(rune(functionName[0])) {
		panic(fmt.Errorf("can not use unexported method %q as second argument to NewFuncReferenceType", functionName))
	}
	return &FuncReferenceType{typeID, functionName}
}

func (t *FuncReferenceType) Arguments() []interface{} {
	return []interface{}{"@" + t.TypeID}
}

func (t *FuncReferenceType) Generate(resolver *ParameterResolver) interface{} {
	referencedType := resolver.Container.Get(t.TypeID)

	v := reflect.ValueOf(referencedType)
	method := v.MethodByName(t.FunctionName)

	if method.IsValid() == false {
		panic(fmt.Errorf("could not generate func reference type @%s::%s method does not exist", t.TypeID, t.FunctionName))
	}

	return method.Interface()
}
