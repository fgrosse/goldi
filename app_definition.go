package goldi

import "fmt"

type AppDefinition map[string]*TypeGenerator

func NewAppDefinition() AppDefinition {
	return AppDefinition{}
}

func (d AppDefinition) RegisterType(typeID string, generatorFunction interface{}) error {
	generator, err := NewTypeGenerator(generatorFunction)
	if err != nil {
		return fmt.Errorf("could register type %q : %s", typeID, err.Error())
	}

	d[typeID] = generator
	return nil
}
