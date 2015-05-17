package goldi

// A TypeFactory is used to instantiate a certain type.
// Its primary implementation is currently the Type type
type TypeFactory interface {

	// Arguments returns all arguments that are used to generate the type.
	// This enables the container validator to check if all required parameters exist
	// and if there are circular type dependencies.
	Arguments() []interface{}

	// Generate will instantiate a new instance of the according type.
	// The given configuration is used to resolve parameters.
	// The type registry can be used to lazily resolve type references.
	Generate(config map[string]interface{}, registry TypeRegistry) interface{}
}
