package goldi

import "fmt"

type Container struct {
	typeRegistry TypeRegistry
	config       map[string]interface{}
	typeCache    map[string]interface{}
}

func NewContainer(registry TypeRegistry, config map[string]interface{}) *Container {
	return &Container{
		typeRegistry: registry,
		config:       config,
		typeCache:    map[string]interface{}{},
	}
}

func (c *Container) Get(typeID string) interface{} {
	t, isCached := c.typeCache[typeID]
	if isCached {
		return t
	}

	generator, isDefined := c.typeRegistry[typeID]
	if isDefined == false {
		panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
	}

	t = generator.Generate(c.config)
	c.typeCache[typeID] = t
	return t
}
