package goldi

import "fmt"

// The ContainerValidator can be used to determine whether a Container is valid.
// A Container is said to be valid if it does not define any type that depends on a
// undefined parameter and does not reference any unregistered type.
// Additionally goldi does not allow you to define circular type references currently.
type ContainerValidator struct {
	checkedTypes               StringSet
	circularDependencyCheckMap StringSet
}

// NewContainerValidator creates a new ContainerValidator
func NewContainerValidator() *ContainerValidator {
	return &ContainerValidator{}
}

// Validate checks if the given container contains any type that fails any of the following checks:
// * it uses a parameter that has not been defined
// * it references a type that has not been defined
// * there is a circular dependency to other types (FooType requires BarType requires BazType requires FooType to be built)
func (v *ContainerValidator) Validate(container *Container) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("container validation failed: %s", err)
		}
	}()

	for typeID, typeDef := range container.TypeRegistry {
		// reset the validation type cache
		v.checkedTypes = StringSet{}

		if err = v.validateTypeParameters(typeID, typeDef, container); err != nil {
			return err
		}

		if err = v.validateTypeReferences(typeID, typeDef, container); err != nil {
			return err
		}
	}
	return nil
}

func (v *ContainerValidator) validateTypeParameters(typeID string, typeDef *Type, container *Container) error {
	typeParameters := typeDef.parameterArguments()
	for _, parameterName := range typeParameters {
		_, isParameterDefined := container.config[parameterName]
		if isParameterDefined == false {
			return fmt.Errorf(`the parameter "%%%s%%" is required by type %q but has not been defined`, parameterName, typeID)
		}
	}
	return nil
}

func (v *ContainerValidator) validateTypeReferences(typeID string, typeDef *Type, container *Container) error {
	typeRefParameters := typeDef.typeReferenceArguments()
	for _, referencedTypeID := range typeRefParameters {
		if v.checkedTypes.Contains(referencedTypeID) {
			continue
		}

		referencedType, err := v.checkTypeIsDefined(typeID, referencedTypeID, container)
		if err != nil {
			return err
		}

		v.circularDependencyCheckMap = StringSet{}
		v.circularDependencyCheckMap.Set(typeID)
		if err = v.checkCircularDependency(referencedType, referencedTypeID, container); err != nil {
			return err
		}

		v.checkedTypes.Set(referencedTypeID)
	}
	return nil
}

func (v *ContainerValidator) checkTypeIsDefined(typeID, referencedTypeID string, container *Container) (*Type, error) {
	typeDef, isDefined := container.TypeRegistry[referencedTypeID]
	if isDefined == false {
		return nil, fmt.Errorf("type %q references unkown type %q", typeID, referencedTypeID)
	}

	return typeDef, nil
}

func (v *ContainerValidator) checkCircularDependency(t *Type, typeID string, container *Container) error {
	typeRefParameters := t.typeReferenceArguments()
	for _, referencedTypeID := range typeRefParameters {
		referencedType, err := v.checkTypeIsDefined(typeID, referencedTypeID, container)
		if err != nil {
			return nil
		}

		if v.circularDependencyCheckMap.Contains(referencedTypeID) {
			return fmt.Errorf("detected circular dependency for type %q (referenced by %q)", referencedTypeID, typeID)
		}

		v.circularDependencyCheckMap.Set(typeID)
		if err = v.checkCircularDependency(referencedType, referencedTypeID, container); err != nil {
			return err
		}
	}

	return nil
}