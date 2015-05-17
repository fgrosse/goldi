package goldi

import "fmt"

// Container is the dependency injection container that can be used by your application to define and get types.
//
// Basically this is just a TypeRegistry with access to the application configuration and the knowledge
// of how to build individual services. Additionally this implements the laziness of the DI using a simple in memory type cache
type Container struct {
	TypeRegistry
	config            map[string]interface{}
	parameterResolver *ParameterResolver
	typeCache         map[string]interface{}
}

// NewContainer creates a new container instance using the provided arguments
func NewContainer(registry TypeRegistry, config map[string]interface{}) *Container {
	return &Container{
		TypeRegistry:      registry,
		config:            config,
		parameterResolver: NewParameterResolver(config, registry),
		typeCache:         map[string]interface{}{},
	}
}

// Get retrieves a previously defined type.
// If the requested typeID is unknown (has not been registered before) Get will panic with an error.
// Since Get can only return interface{} you need to add a type assertion after the call:
//
// 	container.Get("logger").(LoggerInterface)
//
// For your dependency injection to work properly it is important that you do only try to assert interface types
// when you use Get(..). Otherwise it might be impossible to assert the correct type when you change the underlying type
// implementations. Also make sure your application is properly tested and defers some panic handling in case you
// forgot to define a service.
func (c *Container) Get(typeID string) interface{} {
	t, isCached := c.typeCache[typeID]
	if isCached {
		return t
	}

	generator, isDefined := c.TypeRegistry[typeID]
	if isDefined == false {
		panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
	}

	t = generator.Generate(c.parameterResolver)
	c.typeCache[typeID] = t
	return t
}
