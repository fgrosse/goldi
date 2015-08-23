package goldi

import (
	"fmt"
	"reflect"
	"unicode"
)

type funcReferenceType struct {
	typeID string
	functionName string
}

func NewFuncReferenceType(typeID, functionName string) TypeFactory {
	if functionName == "" || unicode.IsLower(rune(functionName[0])) {
		return newInvalidType(fmt.Errorf("can not use unexported method %q as second argument to NewFuncReferenceType", functionName))
	}

	return &funcReferenceType{typeID, functionName}
}

func (t *funcReferenceType) Arguments() []interface{} {
	return []interface{}{"@" + t.typeID}
}

func (t *funcReferenceType) Generate(resolver *ParameterResolver) (interface{}, error) {
	referencedType := resolver.Container.Get(t.typeID)

	v := reflect.ValueOf(referencedType)
	method := v.MethodByName(t.functionName)

	if method.IsValid() == false {
		return nil, fmt.Errorf("could not generate func reference type @%s::%s method does not exist", t.typeID, t.functionName)
	}

	return method.Interface(), nil
}
