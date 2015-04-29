package goldi

import "fmt"

type Container struct {
	definition AppDefinition
	config     map[string]interface{}
}

func NewContainer(definition AppDefinition, config map[string]interface{}) *Container {
	return &Container{definition, config}
}

func (c *Container) Get(typeID string) interface{} {
	generator, isDefined := c.definition[typeID]
	if isDefined {
		return generator.Generate(c.config)
	}

	panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
}
