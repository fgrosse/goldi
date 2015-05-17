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

	if factoryType.NumIn() != len(parameters) {
		panic(fmt.Errorf("invalid number of input parameters: got %d but expected %d", factoryType.NumIn(), len(parameters)))
	}

	return &Type{
		factory:          reflect.ValueOf(function),
		factoryType:      factoryType,
		factoryArguments: buildFactoryCallArguments(factoryType, parameters),
	}
}

func buildFactoryCallArguments(factoryType reflect.Type, factoryParameters []interface{}) []reflect.Value {
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

// Generate will instantiate a new instance of the according type.
// The given configuration is used to resolve parameters that are used in the type factory method
// The type registry is used to lazily resolve type references
func (t *Type) Generate(config map[string]interface{}, registry TypeRegistry) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not generate type: %v", r))
		}
	}()

	if t.factory.IsValid() == false {
		panic("this type is not initialized. Did you use NewType to create it?")
	}

	args := make([]reflect.Value, len(t.factoryArguments))
	for i, argument := range t.factoryArguments {
		args[i] = t.resolveParameter(i, argument, t.factoryType.In(i), config, registry)
	}

	result := t.factory.Call(args)
	if len(result) == 0 {
		// in theory this condition can never evaluate to true since we check the number of return arguments in NewType
		panic(fmt.Errorf("no return parameter found. this should never ever happen ò.Ó"))
	}

	return result[0].Interface()
}

func (t *Type) generateFromFactory(config map[string]interface{}, registry TypeRegistry) interface{} {
	args := make([]reflect.Value, len(t.factoryArguments))
	for i, argument := range t.factoryArguments {
		args[i] = t.resolveParameter(i, argument, t.factoryType.In(i), config, registry)
	}

	result := t.factory.Call(args)
	if len(result) == 0 {
		// in theory this condition can never evaluate to true since we check the number of return arguments in NewType
		panic(fmt.Errorf("no return parameter found. this should never ever happen ò.Ó"))
	}

	return result[0].Interface()
}

func (t *Type) resolveParameter(i int, argument reflect.Value, expectedArgument reflect.Type, config map[string]interface{}, registry TypeRegistry) reflect.Value {
	if argument.Kind() != reflect.String {
		return argument
	}

	stringArgument := argument.Interface().(string)
	if isParameterOrTypeReference(stringArgument) == false {
		return argument
	}

	if stringArgument[0] == '@' {
		return t.resolveTypeReference(i, stringArgument[1:], config, registry, expectedArgument)
	}

	parameterName := stringArgument[1 : len(stringArgument)-1]
	configuredValue, isConfigured := config[parameterName]
	if isConfigured == false {
		return argument
	}

	argument = reflect.New(expectedArgument).Elem()
	argument.Set(reflect.ValueOf(configuredValue))
	return argument
}

func (t *Type) resolveTypeReference(i int, typeID string, config map[string]interface{}, registry TypeRegistry, expectedArgument reflect.Type) reflect.Value {
	referencedType, typeDefined := registry[typeID]
	if typeDefined == false {
		panic(fmt.Errorf("the referenced type \"@%s\" has not been defined", typeID))
	}

	typeInstance := referencedType.Generate(config, registry)
	if reflect.TypeOf(typeInstance).AssignableTo(expectedArgument) == false {
		panic(t.invalidReferencedTypeErr(typeID, typeInstance, i))
	}

	argument := reflect.New(expectedArgument).Elem()
	argument.Set(reflect.ValueOf(typeInstance))
	return argument
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

// typeReferenceArguments is an internal function that returns all factory arguments that are type references
func (t *Type) typeReferenceArguments() []string {
	var typeRefParameters []string
	for _, argument := range t.factoryArguments {
		stringArgument := argument.Interface().(string)
		if isTypeReference(stringArgument) {
			typeRefParameters = append(typeRefParameters, stringArgument[1:])
		}
	}
	return typeRefParameters
}

// parameterArguments is an internal function that returns all factory arguments that are parameters
func (t *Type) parameterArguments() []string {
	var parameterArguments []string
	for _, argument := range t.factoryArguments {
		stringArgument := argument.Interface().(string)
		if isParameter(stringArgument) {
			parameterArguments = append(parameterArguments, stringArgument[1:len(stringArgument)-1])
		}
	}
	return parameterArguments
}

func isParameterOrTypeReference(p string) bool {
	return isParameter(p) || isTypeReference(p)
}

func isParameter(p string) bool {
	if len(p) < 2 {
		return false
	}

	return p[0] == '%' && p[len(p)-1] == '%'
}

func isTypeReference(p string) bool {
	if len(p) < 2 {
		return false
	}

	return p[0] == '@'
}
