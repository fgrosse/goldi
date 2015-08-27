package generator

import (
	"fmt"
	"strings"
)

// RegistrationCode returns the go code that is necessary to register this type
// To avoid any unexpected behavior you should call TypeDefinition.Validate first
func RegistrationCode(t TypeDefinition, typeID, outputPackageName string) string {
	var typeFactoryCode string

	switch {
	case t.FuncName != "" && t.FuncName[0] != '@':
		typeFactoryCode = funcTypeCode(t, outputPackageName)
	case t.FuncName != "" && t.FuncName[0] == '@':
		typeFactoryCode = funcReferenceTypeCode(t)
	case t.AliasForType != "":
		typeFactoryCode = aliasTypeCode(t)
	case t.FactoryMethod != "" && t.FactoryMethod[0] == '@' && strings.Contains(t.FactoryMethod, "::"):
		typeFactoryCode = proxyTypeCode(t)
	case t.FactoryMethod != "":
		typeFactoryCode = factoryTypeCode(t, outputPackageName)
	case t.TypeName != "":
		typeFactoryCode = structTypeCode(t, outputPackageName)
	default:
		panic(fmt.Errorf("can not generate registration code for %+v", t))
	}

	if len(t.Configurator) == 2 {
		configuratorID := t.Configurator[0][1:]
		configuratorMethod := t.Configurator[1]

		typeFactoryCode = fmt.Sprintf("goldi.NewConfiguredType(\n\t\t%s,\n\t\t%q, %q,\n\t)", typeFactoryCode, configuratorID, configuratorMethod)
	}

	return fmt.Sprintf("types.Register(%q, %s)", typeID, typeFactoryCode)
}

func funcTypeCode(t TypeDefinition, outputPackageName string) string {
	funcName := t.FuncName
	if t.Package != outputPackageName {
		funcName = fmt.Sprintf("%s.%s", t.PackageName(), funcName)
	}

	return fmt.Sprintf("goldi.NewFuncType(%s)", funcName)
}

func funcReferenceTypeCode(t TypeDefinition) string {
	parts := strings.SplitN(t.FuncName, "::", 2)
	return fmt.Sprintf("goldi.NewFuncReferenceType(%q, %q)", parts[0][1:], parts[1])
}

func aliasTypeCode(t TypeDefinition) string {
	alias := t.AliasForType
	if alias[0] == '@' {
		alias = alias[1:]
	}
	return fmt.Sprintf("goldi.NewAliasType(%q)", alias)
}

func factoryTypeCode(t TypeDefinition, outputPackageName string) string {
	factoryMethod := t.FactoryMethod
	if t.Package != outputPackageName {
		factoryMethod = fmt.Sprintf("%s.%s", t.PackageName(), t.FactoryMethod)
	}

	arguments := []string{factoryMethod}
	arguments = append(arguments, t.Arguments()...)
	return fmt.Sprintf("goldi.NewType(%s)", strings.Join(arguments, ", "))
}

func structTypeCode(t TypeDefinition, outputPackageName string) string {
	factoryMethod := fmt.Sprintf("new(%s)", t.TypeName)
	if t.Package != outputPackageName {
		factoryMethod = fmt.Sprintf("new(%s.%s)", t.PackageName(), t.TypeName)
	}

	arguments := []string{factoryMethod}
	arguments = append(arguments, t.Arguments()...)
	return fmt.Sprintf("goldi.NewStructType(%s)", strings.Join(arguments, ", "))
}

func proxyTypeCode(t TypeDefinition) string {
	factory := t.FactoryMethod[1:] // omit leading @
	parts := strings.Split(factory, "::")
	arguments := append([]string{fmt.Sprintf("%q", parts[0]), fmt.Sprintf("%q", parts[1])}, t.Arguments()...)
	return fmt.Sprintf("goldi.NewProxyType(%s)", strings.Join(arguments, ", "))
}
