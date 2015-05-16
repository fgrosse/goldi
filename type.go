package goldi

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
)

// A Type holds all information that is necessary to create a new instance of a type ID
type Type struct {
	factory          reflect.Value
	factoryType      reflect.Type
	factoryArguments []reflect.Value
	isFactoryMethod  bool
}

// NewType creates a new Type and checks if the given factory method can be used get a go type
//
// This function will panic if the factoryFunction is no function, returns zero or more than
// one parameter or the return parameter is no pointer or interface type.
// If the number of given factoryParameters does not match the number of arguments of the
// factoryFunction this function will panic as well
func NewType(factoryFunction interface{}, factoryParameters ...interface{}) *Type {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not register type: %v", r))
		}
	}()

	factoryType := reflect.TypeOf(factoryFunction)
	kind := factoryType.Kind()
	switch {
	case kind == reflect.Ptr && factoryType.Elem().Kind() == reflect.Struct:
		return newTypeFromStruct(factoryFunction, factoryType, factoryParameters)
	case kind == reflect.Func:
		return newTypeFromFactoryFunction(factoryFunction, factoryType, factoryParameters)
	default:
		panic(fmt.Errorf("the given factory function must either be a function or a pointer to a struct (given %q)", factoryType.Kind()))
	}
}

func newTypeFromStruct(structFactory interface{}, generatedType reflect.Type, parameters []interface{}) *Type {
	args := make([]reflect.Value, len(parameters))
	for i, argument := range parameters {
		args[i] = reflect.ValueOf(argument)
	}

	return &Type{
		factory:          reflect.ValueOf(structFactory),
		isFactoryMethod:  false,
		factoryType:      generatedType,
		factoryArguments: args,
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
		isFactoryMethod:  true,
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

	if t.isFactoryMethod {
		return t.generateFromFactory(config, registry)
	} else {
		return t.generateFromStruct(config, registry)
	}
}

func (t *Type) generateFromFactory(config map[string]interface{}, registry TypeRegistry) interface{} {
	args := make([]reflect.Value, len(t.factoryArguments))
	for i, argument := range t.factoryArguments {
		args[i] = t.resolveParameter(i, argument, t.factoryType.In(i), config, registry)
	}

	result := t.factory.Call(args)
	if len(result) == 0 {
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

	n := t.factory.Type().NumIn()
	factoryArguments := make([]string, n)
	for i := 0; i < n; i++ {
		arg := t.factory.Type().In(i)
		factoryArguments[i] = arg.String()
	}

	err := fmt.Errorf("the referenced type \"@%s\" (type %T) can not be passed as argument %d to the function signature %s(%s)",
		typeID, typeInstance, i+1, factoryName, strings.Join(factoryArguments, ", "),
	)

	return err
}

func (t *Type) generateFromStruct(config map[string]interface{}, registry TypeRegistry) interface{} {
	if t.factory.Elem().NumField() < len(t.factoryArguments) {
		panic(fmt.Errorf("the struct %T has only %d fields but %d arguments where provided on type registration",
			t.factory.Elem().Interface(), t.factory.Elem().NumField(), len(t.factoryArguments),
		))
	}

	args := make([]reflect.Value, len(t.factoryArguments))
	for i, argument := range t.factoryArguments {
		expectedArgument := t.factory.Elem().Field(i).Type()
		args[i] = t.resolveParameter(i, argument, expectedArgument, config, registry)
	}

	factory := reflect.New(t.factory.Elem().Type())
	n := int(math.Min(float64(factory.Elem().NumField()), float64(len(args))))
	for i := 0; i < n; i++ {
		factory.Elem().Field(i).Set(args[i])
	}

	return factory.Interface()
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
