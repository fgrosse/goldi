package goldi

import "fmt"

type AppDefinition map[string]*Type

func NewAppDefinition() AppDefinition {
	return AppDefinition{}
}

func (d AppDefinition) RegisterType(typeID string, generatorFunction interface{}) error {
	t, err := NewType(generatorFunction)
	if err != nil {
		return err
	}

	return d.Register(typeID, t)
}

func (d AppDefinition) Register(typeID string, typeDef *Type) (err error) {
	_, typeHasAlreadyBeenRegistered := d[typeID]
	if typeHasAlreadyBeenRegistered {
		return fmt.Errorf("type %q has already been registered", typeID)
	}

	d[typeID] = typeDef
	return nil
}
