package goldi

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// A Type holds all information that is necessary to create a new instance of a type.
// Type implements the TypeFactory interface.
type Type struct {
	factory          reflect.Value
	factoryType      reflect.Type
	factoryArguments []reflect.Value
}

// NewType creates a new Type.
//
// This function will panic if:
//   - the factoryFunction is no function,
//   - the factoryFunction returns zero or more than one parameter
//   - the factoryFunctions return parameter is no pointer or interface type.
//   - the number of given factoryParameters does not match the number of arguments of the factoryFunction
func NewType(factoryFunction interface{}, factoryParameters ...interface{}) *Type {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not register type: %v", r))
		}
	}()

	factoryType := reflect.TypeOf(factoryFunction)
	kind := factoryType.Kind()
	switch {
	case kind == reflect.Func:
		return newTypeFromFactoryFunction(factoryFunction, factoryType, factoryParameters)
	default:
		panic(fmt.Errorf("the given factoryFunction must be a function (given %q)", factoryType.Kind()))
	}
}

func newTypeFromFactoryFunction(function interface{}, factoryType reflect.Type, parameters []interface{}) *Type {
	if factoryType.NumOut() != 1 {
		panic(fmt.Errorf("invalid number of return parameters: %d", factoryType.NumOut()))
	}

	kindOfGeneratedType := factoryType.Out(0).Kind()
	if kindOfGeneratedType != reflect.Interface && kindOfGeneratedType != reflect.Ptr {
		panic(fmt.Errorf("return parameter is no interface or pointer but a %v", kindOfGeneratedType))
	}

	if factoryType.IsVariadic() == false && factoryType.NumIn() != len(parameters) {
		panic(fmt.Errorf("invalid number of input parameters: got %d but expected %d", len(parameters), factoryType.NumIn()))
	}

	return &Type{
		factory:          reflect.ValueOf(function),
		factoryType:      factoryType,
		factoryArguments: buildFactoryCallArguments(factoryType, parameters),
	}
}

func buildFactoryCallArguments(factoryType reflect.Type, factoryParameters []interface{}) []reflect.Value {
	if factoryType.IsVariadic() {
		factoryParameters = []interface{}{factoryParameters}
	}

	args := make([]reflect.Value, len(factoryParameters))
	for i, argument := range factoryParameters {
		expectedArgumentType := factoryType.In(i)
		args[i] = reflect.ValueOf(argument)
		if args[i].Kind() != expectedArgumentType.Kind() {
			if stringArg, isString := argument.(string); isString && isParameterOrTypeReference(stringArg) == false {
				panic(fmt.Errorf("input argument %d is of type %s but needs to be a %s", i+1, args[i].Kind(), expectedArgumentType.Kind()))
			}
		}
	}

	return args
}

// Arguments returns all factory parameters from NewType
func (t *Type) Arguments() []interface{} {
	args := make([]interface{}, len(t.factoryArguments))
	for i, argument := range t.factoryArguments {
		args[i] = argument.Interface()
	}
	return args
}

// Generate will instantiate a new instance of the according type.
func (t *Type) Generate(parameterResolver *ParameterResolver) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not generate type: %v", r))
		}
	}()

	if t.factory.IsValid() == false {
		panic("this type is not initialized. Did you use NewType to create it?")
	}

	args := t.generateFactoryArguments(parameterResolver)
	result := t.factory.Call(args)
	if len(result) == 0 {
		// in theory this condition can never evaluate to true since we check the number of return arguments in NewType
		panic(fmt.Errorf("no return parameter found. this should never ever happen ò.Ó"))
	}

	return result[0].Interface()
}

func (t *Type) generateFactoryArguments(parameterResolver *ParameterResolver) []reflect.Value {
	args := make([]reflect.Value, len(t.factoryArguments))
	var err error

	for i, argument := range t.factoryArguments {
		args[i], err = parameterResolver.Resolve(argument, t.factoryType.In(i))

		switch errorType := err.(type) {
		case nil:
			continue
		case TypeReferenceError:
			panic(t.invalidReferencedTypeErr(errorType.TypeID, errorType.TypeInstance, i))
		default:
			panic(err)
		}
	}

	return args
}

func (t *Type) invalidReferencedTypeErr(typeID string, typeInstance interface{}, i int) error {
	factoryName := runtime.FuncForPC(t.factory.Pointer()).Name()
	factoryNameParts := strings.Split(factoryName, "/")
	factoryName = factoryNameParts[len(factoryNameParts)-1]

	n := t.factoryType.NumIn()
	factoryArguments := make([]string, n)
	for i := 0; i < n; i++ {
		arg := t.factoryType.In(i)
		factoryArguments[i] = arg.String()
	}

	err := fmt.Errorf("the referenced type \"@%s\" (type %T) can not be passed as argument %d to the function signature %s(%s)",
		typeID, typeInstance, i+1, factoryName, strings.Join(factoryArguments, ", "),
	)

	return err
}
