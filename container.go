package goldi

import "fmt"

// Container is the dependency injection container that can be used by your application to define and get types.
//
// Basically this is just a TypeRegistry with access to the application configuration and the knowledge
// of how to build individual services. Additionally this implements the laziness of the DI using a simple in memory type cache
//
// You must use goldi.NewContainer to get a initialized instance of a Container!
type Container struct {
	TypeRegistry
	config            map[string]interface{}
	parameterResolver *ParameterResolver
	typeCache         map[string]interface{}
}

// NewContainer creates a new container instance using the provided arguments
func NewContainer(registry TypeRegistry, config map[string]interface{}) *Container {
	c := &Container{
		TypeRegistry: registry,
		config:       config,
		typeCache:    map[string]interface{}{},
	}

	c.parameterResolver = NewParameterResolver(c)
	return c
}

// Get retrieves a previously defined type.
// If the requested typeID is unknown (has not been registered before) Get will panic with an error.
// Since Get can only return interface{} you need to add a type assertion after the call:
//
//     container.Get("logger").(LoggerInterface)
//
// For your dependency injection to work properly it is important that you do only try to assert interface types
// when you use Get(..). Otherwise it might be impossible to assert the correct type when you change the underlying type
// implementations. Also make sure your application is properly tested and defers some panic handling in case you
// forgot to define a service.
func (c *Container) Get(typeID string) interface{} {
	instance, isDefined := c.get(typeID)
	if isDefined == false {
		panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
	}

	return instance
}

func (c *Container) get(typeID string) (interface{}, bool) {
	t, isCached := c.typeCache[typeID]
	if isCached {
		return t, true
	}

	generator, isDefined := c.TypeRegistry[typeID]
	if isDefined == false {
		return nil, false
	}

	instance, err := generator.Generate(c.parameterResolver)
	if err != nil {
		panic(fmt.Errorf("goldi: error while genereating type %q: %s", typeID, err))
	}
	c.typeCache[typeID] = instance
	return instance, true
}
