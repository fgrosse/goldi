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
		// TODO: check argument types
		args[i] = reflect.ValueOf(argument)
	}

	return &StructType{
		structType:   generatedType,
		structFields: args,
	}
}

// Arguments returns all struct parameters from NewStructType
func (t *StructType) Arguments() []interface{} {
	args := make([]interface{}, len(t.structFields))
	for i, argument := range t.structFields {
		args[i] = argument.Interface()
	}
	return args
}

// Generate will instantiate a new instance of the according type.
func (t *StructType) Generate(parameterResolver *ParameterResolver) (interface{}, error) {
	if t.structType == nil {
		panic("this struct type is not initialized. Did you use NewStructType to create it?")
	}

	args, err := t.generateTypeFields(parameterResolver)
	if err != nil {
		return nil, err
	}

	newStructInstance := reflect.New(t.structType)
	for i := 0; i < len(args); i++ {
		newStructInstance.Elem().Field(i).Set(args[i])
	}

	return newStructInstance.Interface(), nil
}

func (t *StructType) generateTypeFields(parameterResolver *ParameterResolver) ([]reflect.Value, error) {
	args := make([]reflect.Value, len(t.structFields))
	var err error

	for i, argument := range t.structFields {
		expectedArgument := t.structType.Field(i).Type
		args[i], err = parameterResolver.Resolve(argument, expectedArgument)

		switch errorType := err.(type) {
		case nil:
			continue
		case TypeReferenceError:
			return nil, t.invalidReferencedTypeErr(errorType.TypeID, errorType.TypeInstance, i)
		default:
			return nil, err
		}
	}

	return args, nil
}

func (t *StructType) invalidReferencedTypeErr(typeID string, typeInstance interface{}, i int) error {
	err := fmt.Errorf("the referenced type \"@%s\" (type %T) can not be used as field %d for struct type %v",
		typeID, typeInstance, i+1, t.structType,
	)

	return err
}
