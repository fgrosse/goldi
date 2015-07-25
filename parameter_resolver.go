package goldi

import "reflect"

// The ParameterResolver is used by type factories to resolve the values of the dynamic factory arguments
// (parameters and other type references).
type ParameterResolver struct {
	Config   map[string]interface{}
	Registry TypeRegistry
}

// NewParameterResolver creates a new ParameterResolver and initializes it with the given config and TypeRegistry.
// The config is used when resolving parameters and the registry for type references.
func NewParameterResolver(config map[string]interface{}, registry TypeRegistry) *ParameterResolver {
	return &ParameterResolver{
		Config:   config,
		Registry: registry,
	}
}

// Resolve takes a parameter and resolves any references to configuration parameter values or type references.
// If the type of `parameter` is not a parameter or type reference it is returned as is.
// Parameters must always have the form `%my.beautiful.param%.
// Type references must have the form `@my_type.bla`
func (r *ParameterResolver) Resolve(parameter reflect.Value, expectedType reflect.Type) (reflect.Value, error) {
	if parameter.Kind() != reflect.String {
		return parameter, nil
	}

	stringParameter := parameter.Interface().(string)
	if isParameterOrTypeReference(stringParameter) == false {
		return parameter, nil
	}

	if isTypeReference(stringParameter) {
		return r.resolveTypeReference(stringParameter[1:], expectedType)
	} else {
		return r.resolveParameter(parameter, stringParameter, expectedType), nil
	}
}

func (r *ParameterResolver) resolveParameter(parameter reflect.Value, stringParameter string, expectedType reflect.Type) reflect.Value {
	parameterName := stringParameter[1 : len(stringParameter)-1]
	configuredValue, isConfigured := r.Config[parameterName]
	if isConfigured == false {
		return parameter
	}

	parameter = reflect.New(expectedType).Elem()
	parameter.Set(reflect.ValueOf(configuredValue))
	return parameter
}

func (r *ParameterResolver) resolveTypeReference(typeID string, expectedType reflect.Type) (reflect.Value, error) {
	referencedType, typeDefined := r.Registry[typeID]
	if typeDefined == false {
		return reflect.Value{}, NewUnknownTypeReferenceError(typeID, `the referenced type "@%s" has not been defined`, typeID)
	}

	typeInstance := referencedType.Generate(r)
	if reflect.TypeOf(typeInstance).AssignableTo(expectedType) == false {
		return reflect.Value{}, NewTypeReferenceError(typeID, typeInstance,
			`the referenced type "@%s" (type %T) is not assignable to the expected type %v`, typeID, typeInstance, expectedType,
		)
	}

	argument := reflect.New(expectedType).Elem()
	argument.Set(reflect.ValueOf(typeInstance))
	return argument, nil
}

func isParameterOrTypeReference(p string) bool {
	return isParameter(p) || isTypeReference(p)
}

func isParameter(p string) bool {
	if len(p) < 3 {
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
