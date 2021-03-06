package goldi

// A TypeFactory is used to instantiate a certain type.
type TypeFactory interface {

	// Arguments returns all arguments that are used to generate the type.
	// This enables the container validator to check if all required parameters exist
	// and if there are circular type dependencies.
	Arguments() []interface{}

	// Generate will instantiate a new instance of the according type or return an error.
	Generate(parameterResolver *ParameterResolver) (interface{}, error)
}
