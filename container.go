package goldi

import "fmt"

type Container struct {
	definition AppDefinition
}

func NewContainer(definition AppDefinition) *Container {
	return &Container{definition}
}

func (c *Container) Get(typeID string) interface{} {
	generator, isDefined := c.definition[typeID]
	if isDefined {
		return generator.Generate()
	}

	panic(fmt.Errorf("could not get type %q : no such type has been defined", typeID))
}
