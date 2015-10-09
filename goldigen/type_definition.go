package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// A TypeDefinition holds all information necessary to register a type for a specific type ID
type TypeDefinition struct {
	Package       string   `yaml:"package"`
	TypeName      string   `yaml:"type"`
	FuncName      string   `yaml:"func"`
	FactoryMethod string   `yaml:"factory"`
	AliasForType  string   `yaml:"alias"`
	Configurator  []string `yaml:"configurator"`

	RawArguments      []interface{} `yaml:"arguments,omitempty"`
	RawArgumentsShort []interface{} `yaml:"args,omitempty"`

	// ForcePackageName can be used in case the full package does not correspond to the actual package name
	ForcePackageName string `yaml:"package-name,omitempty"`
}

// Validate checks if this type definition contains all required fields
func (t *TypeDefinition) Validate(typeID string) error {
	if t.AliasForType != "" {
		return t.validateTypeAlias(typeID)
	}

	if t.FuncName == "" || t.FuncName[0] != '@' {
		if !(t.FactoryMethod != "" && t.FactoryMethod[0] == '@' && strings.Contains(t.FactoryMethod, "::")) {
			if err := t.requireField("package", t.Package, typeID); err != nil {
				return err
			}
		}
	}

	if t.TypeName == "" && t.FuncName == "" {
		if err := t.requireField("factory", t.FactoryMethod, typeID); err != nil {
			return err
		}
	}

	if t.FuncName != "" {
		if t.FactoryMethod != "" {
			return fmt.Errorf("type definition of %q can not have both a factory and a function. Please decide for one of them", typeID)
		}

		if len(t.RawArguments) != 0 {
			return fmt.Errorf("type definition of %q is a function type but contains arguments. Function types do not accept arguments", typeID)
		}
	}

	if len(t.Configurator) > 0 {
		if len(t.Configurator) != 2 {
			return fmt.Errorf("configurator of type %q needs exactly 2 arguments but got %d", typeID, len(t.Configurator))
		}

		if strings.TrimSpace(t.Configurator[0]) == "" || strings.TrimSpace(t.Configurator[1]) == "" {
			return fmt.Errorf("configurator of type %q can not have empty arguments", typeID)
		}

		if t.Configurator[0][0] != '@' {
			return fmt.Errorf("configurator of type %q is no valid type ID (does not start with @)", typeID)
		}

		if unicode.IsLower(rune(t.Configurator[1][0])) {
			return fmt.Errorf("configurator method of type %q is not exported (lowercase)", typeID)
		}
	}

	return nil
}

func (t *TypeDefinition) validateTypeAlias(typeID string) error {
	if t.FactoryMethod != "" {
		return fmt.Errorf("type alias %q must not define a factory method", typeID)
	}

	if t.Package != "" {
		return fmt.Errorf("type alias %q must not define a package name", typeID)
	}

	if t.FuncName != "" {
		return fmt.Errorf("type alias %q must not define a func", typeID)
	}

	if len(t.RawArguments) != 0 {
		return fmt.Errorf("type alias %q must not contain arguments", typeID)
	}

	return nil
}

func (t *TypeDefinition) requireField(fieldName, value, typeID string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("type definition of %q is missing the required %q key", typeID, fieldName)
	}
	return nil
}

var versionSuffix = regexp.MustCompile(`\.v\d+$`)

func (t *TypeDefinition) PackageName() string {
	if t.ForcePackageName != "" {
		return t.ForcePackageName
	}

	pkg := versionSuffix.ReplaceAllString(t.Package, "")
	packageParts := strings.Split(pkg, "/")
	return packageParts[len(packageParts)-1]
}

func (t *TypeDefinition) Arguments() []string {
	rawArgs := append(t.RawArguments, t.RawArgumentsShort...)
	arguments := make([]string, len(rawArgs))
	for i, arg := range rawArgs {
		switch a := arg.(type) {
		case string:
			arguments[i] = fmt.Sprintf(`"%s"`, a)
		default:
			arguments[i] = fmt.Sprintf("%v", a)
		}
	}
	return arguments
}
