package goldi

import (
	"fmt"
	"reflect"
)

type GeneratorFunction func() interface{}

type TypeGenerator struct {
	generator     reflect.Value
	generatedType reflect.Type
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

	return &TypeGenerator{reflect.ValueOf(generator), generatedType}, nil
}

func (g *TypeGenerator) Generate() interface{} {
	result := g.generator.Call([]reflect.Value{})
	if len(result) == 0 {
		panic(fmt.Errorf("could not generate type: no return parameter found. Seems like you did not use goldi.NewTypeGenerator to create this TypeGenerator"))
	}

	return result[0].Interface()
}
