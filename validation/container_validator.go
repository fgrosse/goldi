// Goldi container validation
package validation

import (
	"fmt"

	"github.com/fgrosse/goldi"
)

type ValidationConstraint interface {
	Validate(*goldi.Container) error
}

// The ContainerValidator can be used to determine whether a container passes a set of validation constraints.
type ContainerValidator struct {
	Constraints []ValidationConstraint
}

// NewContainerValidator creates a new ContainerValidator.
// The validator will be initialized with the NoInvalidTypesConstraint, TypeParametersConstraint and TypeReferencesConstraint
func NewContainerValidator() *ContainerValidator {
	return &ContainerValidator{
		Constraints: []ValidationConstraint{
			new(NoInvalidTypesConstraint),
			new(TypeParametersConstraint),
			new(TypeReferencesConstraint),
		},
	}
}

// Add another constraint to this validator
func (v *ContainerValidator) Add(constraint ValidationConstraint) {
	v.Constraints = append(v.Constraints, constraint)
}

// MustValidate behaves exactly as ContainerValidator.Validate but panics if an error occurs
func (v *ContainerValidator) MustValidate(container *goldi.Container) {
	if err := v.Validate(container); err != nil {
		panic(err)
	}
}

// Validate checks if the given container passes all constraints that are registered at the ContainerValidator.
func (v *ContainerValidator) Validate(container *goldi.Container) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("container validation failed: %s", err)
		}
	}()

	for _, constraint := range v.Constraints {
		if err := constraint.Validate(container); err != nil {
			return err
		}
	}

	return nil
}
