package generator

import (
	"fmt"
	"regexp"
	"strings"
)

// A TypeDefinition holds all information necessary to register a type for a specific type ID
type TypeDefinition struct {
	Package       string        `yaml:"package"`
	TypeName      string        `yaml:"type"`
	FuncName      string        `yaml:"func"`
	FactoryMethod string        `yaml:"factory"`
	RawArguments  []interface{} `yaml:"arguments,omitempty"`

	// ForcePackageName can be used in case the full package does not correspond to the actual package name
	ForcePackageName string `yaml:"package-name,omitempty"`
}

// RegistrationCode returns the go code that is necessary to register this type
func (t *TypeDefinition) RegistrationCode(typeID, outputPackageName string) string {
	if t.FuncName != "" {
		funcName := t.FuncName
		if t.Package != outputPackageName {
			funcName = fmt.Sprintf("%s.%s", t.PackageName(), funcName)
		}
		return fmt.Sprintf("types.Register(%q, goldi.NewFuncType(%s))", typeID, funcName)
	}

	var factoryMethod string
	if t.FactoryMethod != "" {
		factoryMethod = t.FactoryMethod
		if t.Package != outputPackageName {
			factoryMethod = fmt.Sprintf("%s.%s", t.PackageName(), t.FactoryMethod)
		}
	} else if t.TypeName != "" {
		factoryMethod = fmt.Sprintf("new(%s)", t.TypeName)
		if t.Package != outputPackageName {
			factoryMethod = fmt.Sprintf("new(%s.%s)", t.PackageName(), t.TypeName)
		}
	}

	arguments := []string{factoryMethod}
	arguments = append(arguments, t.Arguments()...)
	return fmt.Sprintf("types.RegisterType(%q, %s)", typeID, strings.Join(arguments, ", "))
}

// Validate checks if this type definition contains all required fields
func (t *TypeDefinition) Validate(typeID string) error {
	if err := t.requireField("package", t.Package, typeID); err != nil {
		return err
	}

	if t.TypeName == "" && t.FuncName == "" {
		if err := t.requireField("factory", t.FactoryMethod, typeID); err != nil {
			return err
		}
	}

	if t.FuncName != "" {
		if t.FactoryMethod != "" {
			return fmt.Errorf("type definition of %q can not have both a factory and a function. Please decide for one of them")
		}

		if len(t.RawArguments) != 0 {
			return fmt.Errorf("type definition of %q is a function type but contains arguments. Function types do not accept arguments!")
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
