package goldi

import (
	"fmt"
	"reflect"
	"unicode"
)

type funcReferenceType struct {
	typeID *TypeID
}

// NewFuncReferenceType returns a TypeFactory that returns a method of another type as method value (function).
func NewFuncReferenceType(typeID, functionName string) TypeFactory {
	if functionName == "" || unicode.IsLower(rune(functionName[0])) {
		return newInvalidType(fmt.Errorf("can not use unexported method %q as second argument to NewFuncReferenceType", functionName))
	}

	return &funcReferenceType{NewTypeID("@"+typeID + "::" + functionName)}
}

func (t *funcReferenceType) Arguments() []interface{} {
	return []interface{}{"@" + t.typeID.ID}
}

func (t *funcReferenceType) Generate(resolver *ParameterResolver) (interface{}, error) {
	referencedType, err := resolver.Container.Get(t.typeID.ID)
	if err != nil {
		return nil, fmt.Errorf("could not generate func reference type %s : type %s does not exist", t.typeID.ID)
	}

	v := reflect.ValueOf(referencedType)
	method := v.MethodByName(t.typeID.FuncReferenceMethod)

	if method.IsValid() == false {
		return nil, fmt.Errorf("could not generate func reference type %s : method does not exist", t.typeID)
	}

	return method.Interface(), nil
}
