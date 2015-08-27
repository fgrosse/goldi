package main

import (
	"fmt"
	"github.com/fgrosse/goldi/util"
	"sort"
)

// The TypesConfiguration is the struct that holds the complete dependency injection configuration
// as parsed from a yaml file
type TypesConfiguration struct {
	Parameters map[string]string         `yaml:"parameters,omitempty"`
	Types      map[string]TypeDefinition `yaml:"types,omitempty"`
}

// Validate checks if all type definitions of this configuration are valid
func (c *TypesConfiguration) Validate() (err error) {
	if len(c.Types) == 0 {
		return fmt.Errorf("no types have been defined: please define at least one type")
	}

	for typeID, typeDef := range c.Types {
		err = typeDef.Validate(typeID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Packages returns an alphabetically ordered list of unique package names that are referenced by this type configuration.
func (c *TypesConfiguration) Packages(additionalPackages ...string) []string {
	packages := additionalPackages
	seenPackages := util.StringSet{}
	for _, additionalPackage := range additionalPackages {
		seenPackages.Set(additionalPackage)
	}

	for _, typeDef := range c.Types {
		if seenPackages.Contains(typeDef.Package) {
			continue
		}

		seenPackages.Set(typeDef.Package)
		packages = append(packages, typeDef.Package)
	}

	sort.Strings(packages)
	return packages
}
