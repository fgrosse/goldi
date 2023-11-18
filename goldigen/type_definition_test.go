package main_test

import (
	"github.com/fgrosse/goldi/goldigen"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeDefinition", func() {
	Describe("Validate", func() {
		It("should not return an error if the definition contains a factory method", func() {
			t := main.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if the definition contains a type name", func() {
			t := main.TypeDefinition{
				Package:  "foo/bar",
				TypeName: "Baz",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if the definition contains a func name", func() {
			t := main.TypeDefinition{
				Package:  "foo/bar",
				FuncName: "DoFoo",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if the definition contains a type alias", func() {
			t := main.TypeDefinition{
				AliasForType: "type",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should return an error if the definition both contains a func name and a factory method", func() {
			t := main.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewFpp",
				FuncName:      "DoFoo",
			}
			Expect(t.Validate("foobar")).To(MatchError(`type definition of "foobar" can not have both a factory and a function. Please decide for one of them`))
		})

		It("should return an error if the definition is for a func type but contains arguments", func() {
			t := main.TypeDefinition{
				Package:      "foo/bar",
				FuncName:     "DoFoo",
				RawArguments: []interface{}{"test", 42},
			}
			Expect(t.Validate("foobar")).To(MatchError(`type definition of "foobar" is a function type but contains arguments. Function types do not accept arguments`))
		})

		It("should return an error if the definition does not contain a factory method or a type or func name", func() {
			t := main.TypeDefinition{
				Package: "foo/bar",
			}
			Expect(t.Validate("foobar")).To(MatchError(`type definition of "foobar" is missing the required "factory" key`))
		})

		It("should return an error if the configurator does not have exactly two arguments", func() {
			t := main.TypeDefinition{
				Package: "foo/bar", TypeName: "Blup",
				Configurator: []string{"@configurator"},
			}
			Expect(t.Validate("foobar")).To(MatchError(`configurator of type "foobar" needs exactly 2 arguments but got 1`))
		})

		It("should return an error if one of the configurator arguments or both are empty", func() {
			t := main.TypeDefinition{Package: "foo/bar", TypeName: "Blup"}
			invalidArguments := [][]string{{"", ""}, {"@foo", ""}, {"", "Blup"}, {"\t", "  \n "}}
			for _, invalid := range invalidArguments {
				t.Configurator = invalid
				Expect(t.Validate("foobar")).To(MatchError(`configurator of type "foobar" can not have empty arguments`))
			}
		})

		It("should return an error if the configurator type ID does not start with `@`", func() {
			t := main.TypeDefinition{
				Package: "foo/bar", TypeName: "Blup",
				Configurator: []string{"configurator", "Configure"},
			}
			Expect(t.Validate("foobar")).To(MatchError(`configurator of type "foobar" is no valid type ID (does not start with @)`))
		})

		It("should return an error if the configurator method is not exported", func() {
			t := main.TypeDefinition{
				Package:      "foo/bar",
				TypeName:     "Blup",
				Configurator: []string{"@configurator", "configure"},
			}
			Expect(t.Validate("foobar")).To(MatchError(`configurator method of type "foobar" is not exported (lowercase)`))
		})

		It("should not return an error if a func reference type does not contain a package name", func() {
			t := main.TypeDefinition{
				FuncName: "@blup::DoStuff",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if a proxy type does not contain a package name", func() {
			t := main.TypeDefinition{
				FactoryMethod: "@blup::DoStuff",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})
	})

	Describe("PackageName", func() {
		It("should return the package name", func() {
			t := main.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
			}
			Expect(t.PackageName()).To(Equal("bar"))
		})

		It("should strip versions at the end of package names", func() {
			t := main.TypeDefinition{
				Package:       "github.com/fgrosse/servo.v1",
				FactoryMethod: "NewFoo",
			}
			Expect(t.PackageName()).To(Equal("servo"))

			t.Package = "github.com/fgrosse/servov1"
			Expect(t.PackageName()).To(Equal("servov1"))
		})
	})

	Describe("Arguments", func() {
		It("should return all parameters such that they can be used in go code directly", func() {
			t := main.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
				RawArguments: []interface{}{
					"Hello World",
					true,
					42,
					3.1415,
					"%some_parameter%",
					"Hello\t\tWorld",
				},
			}

			arguments := t.Arguments()
			Expect(arguments).To(HaveLen(6))
			Expect(arguments[0]).To(Equal(`"Hello World"`))
			Expect(arguments[1]).To(Equal(`true`))
			Expect(arguments[2]).To(Equal(`42`))
			Expect(arguments[3]).To(Equal(`3.1415`))
			Expect(arguments[4]).To(Equal(`"%some_parameter%"`))
			Expect(arguments[5]).To(Equal("\"Hello\t\tWorld\""))
		})

		It("should return all arguments from RawArgumentsShort", func() {
			t := main.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
				RawArgumentsShort: []interface{}{
					"Hello World",
					true,
					42,
					3.1415,
					"%some_parameter%",
					"Hello\t\tWorld",
				},
			}

			arguments := t.Arguments()
			Expect(arguments).To(HaveLen(6))
			Expect(arguments[0]).To(Equal(`"Hello World"`))
			Expect(arguments[1]).To(Equal(`true`))
			Expect(arguments[2]).To(Equal(`42`))
			Expect(arguments[3]).To(Equal(`3.1415`))
			Expect(arguments[4]).To(Equal(`"%some_parameter%"`))
			Expect(arguments[5]).To(Equal("\"Hello\t\tWorld\""))
		})
	})
})
