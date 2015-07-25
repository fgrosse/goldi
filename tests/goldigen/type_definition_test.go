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

		It("should return an error if the definition contains neither a factory method nor a type name", func() {
			t := generator.TypeDefinition{
				Package: "foo/bar",
			}
			Expect(t.Validate("foobar")).NotTo(Succeed())
		})
	})

	It("should return the package name", func() {
		t := generator.TypeDefinition{
			Package:       "foo/bar",
			FactoryMethod: "NewBaz",
		}
		Expect(t.PackageName()).To(Equal("bar"))
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

	Describe("Factory", func() {
		It("should return the factory function", func() {
			t := generator.TypeDefinition{
				Package:       "foo/bar",
				FactoryMethod: "NewBaz",
			}
			Expect(t.Factory("some/package/lib")).To(Equal("bar.NewBaz"))
			Expect(t.Factory("foo/bar")).To(Equal("NewBaz"))
		})

		It("should return the type struct if no factory function is given", func() {
			t := generator.TypeDefinition{
				Package:  "foo/bar",
				TypeName: "Baz",
			}
			Expect(t.Factory("some/package/lib")).To(Equal("new(bar.Baz)"))
			Expect(t.Factory("foo/bar")).To(Equal("new(Baz)"))
		})
	})
})
