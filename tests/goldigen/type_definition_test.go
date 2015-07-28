package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/generator"
)

var _ = Describe("TypeDefinition", func() {
	Describe("Validate", func() {
		It("should not return an error if the definition contains a factory method", func() {
			t := generator.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if the definition contains a type name", func() {
			t := generator.TypeDefinition{
				Package:  "foo/bar",
				TypeName: "Baz",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should not return an error if the definition contains a func name", func() {
			t := generator.TypeDefinition{
				Package:  "foo/bar",
				FuncName: "DoFoo",
			}
			Expect(t.Validate("foobar")).To(Succeed())
		})

		It("should return an error if the definition both contains a func name and a factory method", func() {
			t := generator.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewFpp",
				FuncName:      "DoFoo",
			}
			Expect(t.Validate("foobar")).NotTo(Succeed())
		})

		It("should return an error if the definition is for a func type but contains arguments", func() {
			t := generator.TypeDefinition{
				Package:       "foo/bar",
				FuncName:      "DoFoo",
				RawArguments: []interface{}{"test", 42},
			}
			Expect(t.Validate("foobar")).NotTo(Succeed())
		})

		It("should return an error if the definition does not contain a factory method or a type or func name", func() {
			t := generator.TypeDefinition{
				Package: "foo/bar",
			}
			Expect(t.Validate("foobar")).NotTo(Succeed())
		})
	})

	Describe("PackageName", func() {
		It("should return the package name", func() {
			t := generator.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
			}
			Expect(t.PackageName()).To(Equal("bar"))
		})

		It("should strip versions at the end of package names", func() {
			t := generator.TypeDefinition{
				Package:       "github.com/fgrosse/servo.v1",
				FactoryMethod: "NewFoo",
			}
			Expect(t.PackageName()).To(Equal("servo"))

			t.Package = "github.com/fgrosse/servov1"
			Expect(t.PackageName()).To(Equal("servov1"))
		})
	})

	It("should return all parameters such that they can be used in go code directly", func() {
		t := generator.TypeDefinition{
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

	Describe("RegistrationCode", func() {
		var typeDef generator.TypeDefinition
		BeforeEach(func() {
			typeDef = generator.TypeDefinition{
				Package:  "foo/bar",
				RawArguments: []interface{}{"foo", "%bar%", 42},
			}
		})

		It("should return the golang code to register a struct type", func() {
			typeDef.TypeName = "Baz"
			Expect(typeDef.RegistrationCode("test_type", "some/package/lib")).To(Equal(`types.RegisterType("test_type", new(bar.Baz), "foo", "%bar%", 42)`))
			Expect(typeDef.RegistrationCode("test_type", typeDef.Package)).To(Equal(`types.RegisterType("test_type", new(Baz), "foo", "%bar%", 42)`))
		})

		It("should return the golang code to register a type using a factory function", func() {
			typeDef.FactoryMethod = "NewBaz"
			Expect(typeDef.RegistrationCode("test_type", "some/package/lib")).To(Equal(`types.RegisterType("test_type", bar.NewBaz, "foo", "%bar%", 42)`))
			Expect(typeDef.RegistrationCode("test_type", typeDef.Package)).To(Equal(`types.RegisterType("test_type", NewBaz, "foo", "%bar%", 42)`))
		})

		It("should return the golang code to register a function type", func() {
			typeDef.FuncName = "DoFoo"
			typeDef.RawArguments = nil
			Expect(typeDef.RegistrationCode("test_type", "some/package/lib")).To(Equal(`types.Register("test_type", goldi.NewFuncType(bar.DoFoo))`))
			Expect(typeDef.RegistrationCode("test_type", typeDef.Package)).To(Equal(`types.Register("test_type", goldi.NewFuncType(DoFoo))`))
		})
	})
})
