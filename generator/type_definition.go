package generator

import (
	"fmt"
	"strings"
)

// A TypeDefinition holds all information necessary to register a type for a specific type ID
type TypeDefinition struct {
	Package       string        `yaml:"package"`
	TypeName      string        `yaml:"type"`
	FactoryMethod string        `yaml:"factory"`
	RawArguments  []interface{} `yaml:"arguments,omitempty"`
}

func (t *TypeDefinition) Factory(outputPackageName string) string {
	var factoryMethod string

	if t.FactoryMethod != "" {
		factoryMethod = t.FactoryMethod
		if t.Package != outputPackageName {
			factoryMethod = fmt.Sprintf("%s.%s", t.PackageName(), t.FactoryMethod)
		}
	} else if t.TypeName != "" {
		factoryMethod = fmt.Sprintf("%s{}", t.TypeName)
		if t.Package != outputPackageName {
			factoryMethod = fmt.Sprintf("%s.%s{}", t.PackageName(), t.TypeName)
		}
	}

	return factoryMethod
}

// Validate checks if this type definition contains all required fields
func (t *TypeDefinition) Validate(typeID string) error {
	if err := t.requireField("package", t.Package, typeID); err != nil {
		return err
	}

	if t.TypeName == "" {
		if err := t.requireField("factory", t.FactoryMethod, typeID); err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeDefinition) requireField(fieldName, value, typeId string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("type definition of %q is missing the required %q key", typeId, fieldName)
	}
	return nil
}

func (t *TypeDefinition) PackageName() string {
	packageParts := strings.Split(t.Package, "/")
	return packageParts[len(packageParts)-1]
}

func (t *TypeDefinition) Arguments() []string {
	arguments := make([]string, len(t.RawArguments))
	for i, arg := range t.RawArguments {
		switch a := arg.(type) {
		case string:
			arguments[i] = fmt.Sprintf(`"%s"`, a)
		default:
			arguments[i] = fmt.Sprintf("%v", a)
		}
	}
	return arguments
}
