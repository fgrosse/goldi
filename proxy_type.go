package goldi

import (
	"fmt"
	"unicode"
	"reflect"
)

type proxyType struct {
	typeID *TypeID
	args []interface{}
}

// NewProxyType returns a TypeFactory that uses a function of another type to generate a result.
//
// Goldigen yaml syntax example:
//     logger:
//         factory: "@logger_provider::GetLogger"
//         args:    [ "My Logger" ]
func NewProxyType(typeID, functionName string, args ...interface{}) TypeFactory {
	if functionName == "" || unicode.IsLower(rune(functionName[0])) {
		return newInvalidType(fmt.Errorf("can not use unexported method %q as second argument to NewProxyType", functionName))
	}

	return &proxyType{
		typeID: NewTypeID("@"+typeID + "::" + functionName),
		args:   args,
	}
}

func (t *proxyType) Arguments() []interface{} {
	args := make([]interface{}, len(t.args)+1)
	args[0] = "@" + t.typeID.ID
	for i, a := range t.args {
		args[i+1] = a
	}
	return args
}

func (t *proxyType) Generate(resolver *ParameterResolver) (interface{}, error) {
	referencedType, err := resolver.Container.Get(t.typeID.ID)
	if err != nil {
		return nil, fmt.Errorf("could not generate proxy type %s : type %s does not exist", t.typeID, t.typeID.ID)
	}

	v := reflect.ValueOf(referencedType)
	method := v.MethodByName(t.typeID.FuncReferenceMethod)

	if method.IsValid() == false {
		return nil, fmt.Errorf("could not generate proxy type %s : method does not exist", t.typeID)
	}

	t2 := NewType(method.Interface(), t.args...)
	return t2.Generate(resolver)
}
