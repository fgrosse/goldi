package goldi

import (
	"fmt"
	"reflect"
)

// A StructType holds all information that is necessary to create a new instance of some struct type.
// StructType implements the TypeFactory interface.
type StructType struct {
	structType   reflect.Type
	structFields []reflect.Value
}

// NewStructType creates a new StructType.
//
// This function will panic if:
//   - structT is no struct or pointer to a struct,
//   - the number of given structParameters exceed the number of field of structT
//   - the structParameters types do not match the fields of structT
func NewStructType(structT interface{}, structParameters ...interface{}) *StructType {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not register struct type: %v", r))
		}
	}()

	structType := reflect.TypeOf(structT)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	switch structType.Kind() {
	case reflect.Struct:
		return newTypeFromStruct(structType, structParameters)
	default:
		panic(fmt.Errorf("the given type must either be a struct or a pointer to a struct (given %T)", structT))
	}
}

func newTypeFromStruct(generatedType reflect.Type, parameters []interface{}) *StructType {
	if generatedType.NumField() < len(parameters) {
		panic(fmt.Errorf("the struct %s has only %d fields but %d arguments where provided",
			generatedType.Name(), generatedType.NumField(), len(parameters),
		))
	}

	args := make([]reflect.Value, len(parameters))
	for i, argument := range parameters {
		// TODO check argument types
		args[i] = reflect.ValueOf(argument)
	}

	return &StructType{
		structType:   generatedType,
		structFields: args,
	}
}

// Generate will instantiate a new instance of the according type.
// The given configuration is used to resolve parameters that are used in the type factory method
// The type registry is used to lazily resolve type references
func (t *StructType) Generate(config map[string]interface{}, registry TypeRegistry) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not generate type: %v", r))
		}
	}()

	if t.structType == nil {
		panic("this struct type is not initialized. Did you use NewStructType to create it?")
	}

	args := t.generateTypeFields(config, registry)

	newStructInstance := reflect.New(t.structType)
	for i := 0; i < len(args); i++ {
		newStructInstance.Elem().Field(i).Set(args[i])
	}

	return newStructInstance.Interface()
}

func (t *StructType) generateTypeFields(config map[string]interface{}, registry TypeRegistry) []reflect.Value {
	args := make([]reflect.Value, len(t.structFields))
	for i, argument := range t.structFields {
		expectedArgument := t.structType.Field(i).Type
		args[i] = t.resolveParameter(i, argument, expectedArgument, config, registry)
	}

	return args
}

// TODO refactor this out into a parameter resolver
func (t *StructType) resolveParameter(i int, argument reflect.Value, expectedArgument reflect.Type, config map[string]interface{}, registry TypeRegistry) reflect.Value {
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

func (t *StructType) resolveTypeReference(i int, typeID string, config map[string]interface{}, registry TypeRegistry, expectedArgument reflect.Type) reflect.Value {
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

func (t *StructType) invalidReferencedTypeErr(typeID string, typeInstance interface{}, i int) error {
	err := fmt.Errorf("the referenced type \"@%s\" (type %T) can not be used as field %d for struct type %v",
		typeID, typeInstance, i+1, t.structType,
	)

	return err
}

// typeReferenceArguments is an internal function that returns all factory arguments that are type references
// TODO the container validator needs this. refactor this code duplication!
func (t *StructType) typeReferenceArguments() []string {
	var typeRefParameters []string
	for _, argument := range t.structFields {
		stringArgument := argument.Interface().(string)
		if isTypeReference(stringArgument) {
			typeRefParameters = append(typeRefParameters, stringArgument[1:])
		}
	}
	return typeRefParameters
}

// parameterArguments is an internal function that returns all factory arguments that are parameters
// TODO the container validator needs this. refactor this code duplication!
func (t *StructType) parameterArguments() []string {
	var parameterArguments []string
	for _, argument := range t.structFields {
		stringArgument := argument.Interface().(string)
		if isParameter(stringArgument) {
			parameterArguments = append(parameterArguments, stringArgument[1:len(stringArgument)-1])
		}
	}
	return parameterArguments
}
