package goldi

import "fmt"

type TypeRegistry map[string]*Type

func NewTypeRegistry() TypeRegistry {
	return TypeRegistry{}
}

func (r TypeRegistry) RegisterType(typeID string, generatorFunction interface{}, arguments ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not register type: %s", r)
		}
	}()

	t := NewType(generatorFunction, arguments...)
	r.Register(typeID, t)
	return nil
}

func (r TypeRegistry) Register(typeID string, typeDef *Type) {
	r[typeID] = typeDef
}
