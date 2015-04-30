package goldi

import (
	"fmt"
	"reflect"
)

type GeneratorFunction func() interface{}

type Type struct {
	generator          reflect.Value
	generatorType      reflect.Type
	generatorArguments []reflect.Value
}

func NewType(factoryFunction interface{}, factoryParameters ...interface{}) *Type {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not register type: %v", r))
		}
	}()

	generatorType := reflect.TypeOf(factoryFunction)
	if generatorType.Kind() != reflect.Func {
		panic(fmt.Errorf("kind was %v, not Func", generatorType.Kind()))
	}

	if generatorType.NumOut() != 1 {
		panic(fmt.Errorf("invalid number of return parameters: %d", generatorType.NumOut()))
	}

	kindOfGeneratedType := generatorType.Out(0).Kind()
	if kindOfGeneratedType != reflect.Interface && kindOfGeneratedType != reflect.Ptr {
		panic(fmt.Errorf("return parameter is no interface or pointer but a %v", kindOfGeneratedType))
	}

	if generatorType.NumIn() != len(factoryParameters) {
		panic(fmt.Errorf("invalid number of input parameters: got %d but expected %d", generatorType.NumIn(), len(factoryParameters)))
	}

	return &Type{
		generator:          reflect.ValueOf(factoryFunction),
		generatorType:      generatorType,
		generatorArguments: buildGeneratorCallArguments(generatorType, factoryParameters),
	}
}

func buildGeneratorCallArguments(generatorType reflect.Type, factoryParameters []interface{}) []reflect.Value {
	args := make([]reflect.Value, len(factoryParameters))
	for i, argument := range factoryParameters {
		expectedArgumentType := generatorType.In(i)
		args[i] = reflect.ValueOf(argument)
		if args[i].Kind() != expectedArgumentType.Kind() {
			if stringArg, isString := argument.(string); isString && isParameter(stringArg) == false {
				panic(fmt.Errorf("input argument %d is of type %s but needs to be a %s", i+1, args[i].Kind(), expectedArgumentType.Kind()))
			}
		}
	}

	return args
}

func (t *Type) Generate(config map[string]interface{}) interface{} {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("could not generate type: %v", r))
		}
	}()

	args := make([]reflect.Value, len(t.generatorArguments))
	for i, argument := range t.generatorArguments {
		args[i] = t.resolveParameter(argument, t.generatorType.In(i), config)
	}

	result := t.generator.Call(args)
	if len(result) == 0 {
		panic(fmt.Errorf("no return parameter found. Seems like you did not use goldi.NewTypeGenerator to create this TypeGenerator"))
	}

	return result[0].Interface()
}

func (t *Type) resolveParameter(argument reflect.Value, expectedArgument reflect.Type, config map[string]interface{}) reflect.Value {
	if argument.Kind() != reflect.String {
		return argument
	}

	stringArgument := argument.Interface().(string)
	if isParameter(stringArgument) == false {
		return argument
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
