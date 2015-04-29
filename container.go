package goldi

import "fmt"

type Container struct {
	typeRegistry TypeRegistry
	config       map[string]interface{}
}

func NewContainer(typeRegistry TypeRegistry, config map[string]interface{}) *Container {
	return &Container{typeRegistry, config}
}

func (c *Container) Get(typeID string) interface{} {
	generator, isDefined := c.typeRegistry[typeID]
	if isDefined {
		return generator.Generate(c.config)
	}

	panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
}
