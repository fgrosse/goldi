package goldi

import (
	"fmt"
	"github.com/fgrosse/goldi/util"
)

type ValidationConstraint interface {
	Validate(*Container) error
}

// The ContainerValidator can be used to determine whether a container passes a set of validation constraints.
type ContainerValidator struct {
	constraints []ValidationConstraint
}

// NewContainerValidator creates a new ContainerValidator.
// The validator will be initialized with the TypeParametersConstraint and TypeReferencesConstraint
func NewContainerValidator() *ContainerValidator {
	return &ContainerValidator{
		constraints: []ValidationConstraint{
			new(TypeParametersConstraint),
			new(TypeReferencesConstraint),
		},
	}
}

// Add another constraint to this validator
func (v *ContainerValidator) Add(constraint ValidationConstraint) {
	v.constraints = append(v.constraints, constraint)
}

// MustValidate behaves exactly as ContainerValidator.Validate but panics if an error occurs
func (v *ContainerValidator) MustValidate(container *Container) {
	if err := v.Validate(container); err != nil {
		panic(err)
	}
}

// Validate checks if the given container passes all constraints that are registered at the ContainerValidator.
func (v *ContainerValidator) Validate(container *Container) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("container validation failed: %s", err)
		}
	}()

	for _, constraint := range v.constraints {
		if err := constraint.Validate(container); err != nil {
			return err
		}
	}

	return nil
}

type TypeParametersConstraint struct{}

func (c *TypeParametersConstraint) Validate(container *Container) (err error) {
	for typeID, typeFactory := range container.TypeRegistry {
		allArguments := typeFactory.Arguments()
		if err = c.validateTypeParameters(typeID, container, allArguments); err != nil {
			return err
		}
	}

	return nil
}

func (c *TypeParametersConstraint) validateTypeParameters(typeID string, container *Container, allArguments []interface{}) error {
	typeParameters := c.parameterArguments(allArguments)
	for _, parameterName := range typeParameters {
		_, isParameterDefined := container.config[parameterName]
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
		if isString && isParameter(stringArgument) {
			parameterArguments = append(parameterArguments, stringArgument[1:len(stringArgument)-1])
		}
	}
	return parameterArguments
}

type TypeReferencesConstraint struct {
	checkedTypes               util.StringSet
	circularDependencyCheckMap util.StringSet
}

func (c *TypeReferencesConstraint) Validate(container *Container) (err error) {
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

func (c *TypeReferencesConstraint) validateTypeReferences(typeID string, container *Container, allArguments []interface{}) error {
	typeRefParameters := c.typeReferenceArguments(allArguments)
	for _, referencedTypeID := range typeRefParameters {
		if c.checkedTypes.Contains(referencedTypeID) {
			// TEST: test this for improved code coverage
			continue
		}

		referencedTypeFactory, err := c.checkTypeIsDefined(newTypeId(typeID), newTypeId(referencedTypeID), container)
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
		if isString && isTypeReference(stringArgument) {
			typeRefParameters = append(typeRefParameters, stringArgument[1:])
		}
	}
	return typeRefParameters
}

func (c *TypeReferencesConstraint) checkTypeIsDefined(t, referencedType typeIDT, container *Container) (TypeFactory, error) {
	typeDef, isDefined := container.TypeRegistry[referencedType.ID]
	if isDefined == false {
		return nil, fmt.Errorf("type %q references unknown type %q", t.ID, referencedType.ID)
	}

	return typeDef, nil
}

func (c *TypeReferencesConstraint) checkCircularDependency(typeFactory TypeFactory, typeID string, container *Container) error {
	allArguments := typeFactory.Arguments()
	typeRefParameters := c.typeReferenceArguments(allArguments)

	for _, referencedTypeID := range typeRefParameters {
		referencedType, err := c.checkTypeIsDefined(newTypeId(typeID), newTypeId(referencedTypeID), container)
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
