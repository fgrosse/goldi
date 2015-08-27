package validation

import (
	"fmt"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/util"
)

// The TypeReferencesConstraint is used in a ContainerValidator to check if all referenced types in the container have been defined.
type TypeReferencesConstraint struct {
	checkedTypes               util.StringSet
	circularDependencyCheckMap util.StringSet
}

func (c *TypeReferencesConstraint) Validate(container *goldi.Container) (err error) {
	for typeID, typeFactory := range container.TypeRegistry {
		// reset the validation type cache
		c.checkedTypes = util.StringSet{}
		allArguments := typeFactory.Arguments()

		if err = c.validateTypeReferences(typeID, container, allArguments); err != nil {
			return err
		}
	}

	return nil
}

func (c *TypeReferencesConstraint) validateTypeReferences(typeID string, container *goldi.Container, allArguments []interface{}) error {
	typeRefParameters := c.typeReferenceArguments(allArguments)
	for _, referencedTypeID := range typeRefParameters {
		if c.checkedTypes.Contains(referencedTypeID) {
			// TEST: test this for improved code coverage
			continue
		}

		referencedTypeFactory, err := c.checkTypeIsDefined(goldi.NewTypeID(typeID).ID, goldi.NewTypeID(referencedTypeID).ID, container)
		if err != nil {
			return err
		}

		c.circularDependencyCheckMap = util.StringSet{}
		c.circularDependencyCheckMap.Set(typeID)
		if err = c.checkCircularDependency(referencedTypeFactory, referencedTypeID, container); err != nil {
			return err
		}

		c.checkedTypes.Set(referencedTypeID)
	}
	return nil
}

func (c *TypeReferencesConstraint) typeReferenceArguments(allArguments []interface{}) []string {
	var typeRefParameters []string
	for _, argument := range allArguments {
		stringArgument, isString := argument.(string)
		if isString && goldi.IsTypeReference(stringArgument) {
			typeRefParameters = append(typeRefParameters, stringArgument[1:])
		}
	}
	return typeRefParameters
}

func (c *TypeReferencesConstraint) checkTypeIsDefined(t, referencedType string, container *goldi.Container) (goldi.TypeFactory, error) {
	typeDef, isDefined := container.TypeRegistry[referencedType]
	if isDefined == false {
		return nil, fmt.Errorf("type %q references unknown type %q", t, referencedType)
	}

	return typeDef, nil
}

func (c *TypeReferencesConstraint) checkCircularDependency(typeFactory goldi.TypeFactory, typeID string, container *goldi.Container) error {
	allArguments := typeFactory.Arguments()
	typeRefParameters := c.typeReferenceArguments(allArguments)

	for _, referencedTypeID := range typeRefParameters {
		referencedType, err := c.checkTypeIsDefined(goldi.NewTypeID(typeID).ID, goldi.NewTypeID(referencedTypeID).ID, container)
		if err != nil {
			// TEST: test this for improved code coverage
			return nil
		}

		if c.circularDependencyCheckMap.Contains(referencedTypeID) {
			return fmt.Errorf("detected circular dependency for type %q (referenced by %q)", referencedTypeID, typeID)
		}

		c.circularDependencyCheckMap.Set(typeID)
		if err = c.checkCircularDependency(referencedType, referencedTypeID, container); err != nil {
			return err
		}
	}

	return nil
}
