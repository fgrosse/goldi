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
	Config   map[string]interface{}
	Resolver *ParameterResolver

	typeCache map[string]interface{}
}

// NewContainer creates a new container instance using the provided arguments
func NewContainer(registry TypeRegistry, config map[string]interface{}) *Container {
	c := &Container{
		TypeRegistry: registry,
		Config:       config,
		typeCache:    map[string]interface{}{},
	}

	c.Resolver = NewParameterResolver(c)
	return c
}

// MustGet behaves exactly like Get but will panic instead of returning an error
// Since MustGet can only return interface{} you need to add a type assertion after the call:
//     container.MustGet("logger").(LoggerInterface)
func (c *Container) MustGet(typeID string) interface{} {
	t, err := c.Get(typeID)
	if err != nil {
		panic(err)
	}

	return t
}

// Get retrieves a previously defined type or an error.
// If the requested typeID has not been registered before or can not be generated Get will return an error.
//
// For your dependency injection to work properly it is important that you do only try to assert interface types
// when you use Get(..). Otherwise it might be impossible to assert the correct type when you change the underlying type
// implementations. Also make sure your application is properly tested and defers some panic handling in case you
// forgot to define a service.
//
// See also Container.MustGet
func (c *Container) Get(typeID string) (interface{}, error) {
	instance, isDefined, err := c.get(typeID)
	if err != nil {
		return nil, err
	}

	if isDefined == false {
		return nil, newUnknownTypeReferenceError(typeID, "no such type has been defined")
	}

	return instance, nil
}

func (c *Container) get(typeID string) (interface{}, bool, error) {
	t, isCached := c.typeCache[typeID]
	if isCached {
		return t, true, nil
	}

	generator, isDefined := c.TypeRegistry[typeID]
	if isDefined == false {
		return nil, false, nil
	}

	instance, err := generator.Generate(c.Resolver)
	if err != nil {
		return nil, false, fmt.Errorf("goldi: error while generating type %q: %s", typeID, err)
	}

	c.typeCache[typeID] = instance
	return instance, true, nil
}
