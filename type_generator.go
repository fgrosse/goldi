package goldi

import (
	"fmt"
	"reflect"
)

type GeneratorFunction func() interface{}

type TypeGenerator struct {
	generatorType reflect.Type
	generator     reflect.Value
}

func NewTypeGenerator(generator interface{}) (*TypeGenerator, error) {
	reflectedGenerator := reflect.TypeOf(generator)
	if reflectedGenerator.Kind() != reflect.Func {
		return nil, fmt.Errorf("kind was %v, not Func", reflectedGenerator.Kind())
	}

	if reflectedGenerator.NumOut() != 1 {
		return nil, fmt.Errorf("invalid number of return parameters: %d", reflectedGenerator.NumOut())
	}

	generatedType := reflectedGenerator.Out(0)
	kindOfGeneratedType := generatedType.Kind()
	if kindOfGeneratedType != reflect.Interface && kindOfGeneratedType != reflect.Ptr {
		return nil, fmt.Errorf("return parameter is no interface but a %v", kindOfGeneratedType)
	}

	return &TypeGenerator{reflectedGenerator, reflect.ValueOf(generator)}, nil
}

func (g *TypeGenerator) Generate(args ...interface{}) interface{} {
	defer g.panicHandler()

	arguments := make([]reflect.Value, len(args))
	for i, argument := range args {
		expectedArgumentType := g.generatorType.In(i)
		arguments[i] = reflect.ValueOf(argument)
		if arguments[i].Kind() != expectedArgumentType.Kind() {
			panic(fmt.Errorf("input argument %d is of type %s but needs to be a %s", i+1, arguments[i].Kind(), expectedArgumentType.Kind()))
		}
	}

	result := g.generator.Call(arguments)
	if len(result) == 0 {
		panic(fmt.Errorf("could not generate type: no return parameter found. Seems like you did not use goldi.NewTypeGenerator to create this TypeGenerator"))
	}

	return result[0].Interface()
}

func (g *TypeGenerator) panicHandler() {
	if r := recover(); r != nil {
		panic(fmt.Errorf("goldi type generator: could not generate type: %v", r))
	}
}
