package goldi

import "fmt"

type Container struct {
	TypeRegistry
	config    map[string]interface{}
	typeCache map[string]interface{}
}

func NewContainer(registry TypeRegistry, config map[string]interface{}) *Container {
	return &Container{
		TypeRegistry: registry,
		config:       config,
		typeCache:    map[string]interface{}{},
	}
}

func (c *Container) Get(typeID string) interface{} {
	t, isCached := c.typeCache[typeID]
	if isCached {
		return t
	}

	generator, isDefined := c.TypeRegistry[typeID]
	if isDefined == false {
		panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
	}

	t = generator.Generate(c.config)
	c.typeCache[typeID] = t
	return t
}
