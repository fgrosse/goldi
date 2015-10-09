package validation

import (
	"fmt"

	"github.com/fgrosse/goldi"
)

// The TypeParametersConstraint is used in a ContainerValidator to check if all used parameters do exist.
type TypeParametersConstraint struct{}

// Validate implements the Constraint interface by checking if all referenced parameters have been defined.
func (c *TypeParametersConstraint) Validate(container *goldi.Container) (err error) {
	for typeID, typeFactory := range container.TypeRegistry {
		allArguments := typeFactory.Arguments()
		if err = c.validateTypeParameters(typeID, container, allArguments); err != nil {
			return err
		}
	}

	return nil
}

func (c *TypeParametersConstraint) validateTypeParameters(typeID string, container *goldi.Container, allArguments []interface{}) error {
	typeParameters := c.parameterArguments(allArguments)
	for _, parameterName := range typeParameters {
		_, isParameterDefined := container.Config[parameterName]
		if isParameterDefined == false {
			return fmt.Errorf(`the parameter "%%%s%%" is required by type %q but has not been defined`, parameterName, typeID)
		}
	}
	return nil
}

func (c *TypeParametersConstraint) parameterArguments(allArguments []interface{}) []string {
	var parameterArguments []string
	for _, argument := range allArguments {
		stringArgument, isString := argument.(string)
		if isString && goldi.IsParameter(stringArgument) {
			parameterArguments = append(parameterArguments, stringArgument[1:len(stringArgument)-1])
		}
	}
	return parameterArguments
}
