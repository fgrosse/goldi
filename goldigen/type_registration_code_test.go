package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/goldigen"
)

var _ = Describe("FactoryCode", func() {
	It("should return the golang code to register a struct type", func() {
		typeDef := main.TypeDefinition{
			Package:      "foo/bar",
			TypeName:     "Baz",
			RawArguments: []interface{}{"foo", "%bar%", 42},
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewStructType(new(bar.Baz), "foo", "%bar%", 42)`))
		Expect(main.FactoryCode(typeDef, typeDef.Package)).To(Equal(`goldi.NewStructType(new(Baz), "foo", "%bar%", 42)`))
	})

	It("should return the golang code to register a type using a factory function", func() {
		typeDef := main.TypeDefinition{
			Package:       "foo/bar",
			FactoryMethod: "NewBaz",
			RawArguments:  []interface{}{"foo", "%bar%", 42},
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewType(bar.NewBaz, "foo", "%bar%", 42)`))
		Expect(main.FactoryCode(typeDef, typeDef.Package)).To(Equal(`goldi.NewType(NewBaz, "foo", "%bar%", 42)`))
	})

	It("should return the golang code to register a function type", func() {
		typeDef := main.TypeDefinition{
			Package:  "foo/bar",
			FuncName: "DoFoo",
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewFuncType(bar.DoFoo)`))
		Expect(main.FactoryCode(typeDef, typeDef.Package)).To(Equal(`goldi.NewFuncType(DoFoo)`))
	})

	It("should return the golang code to register a type alias", func() {
		typeDef := main.TypeDefinition{
			AliasForType: "@test_type",
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewAliasType("test_type")`))
	})

	It("should return the golang code to register a func reference type", func() {
		typeDef := main.TypeDefinition{
			FuncName: "@my_controller::FancyAction",
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewFuncReferenceType("my_controller", "FancyAction")`))
	})

	It("should return the golang code to register a proxy type", func() {
		typeDef := main.TypeDefinition{
			FactoryMethod: "@logger_provider::GetLogger",
			RawArguments:  []interface{}{"foo", "%bar%", 42},
		}
		Expect(main.FactoryCode(typeDef, "some/package/lib")).To(Equal(`goldi.NewProxyType("logger_provider", "GetLogger", "foo", "%bar%", 42)`))
	})
})
