package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/goldigen"
)

var _ = Describe("TypesConfiguration", func() {
	Describe("validation", func() {
		It("should return an error if no types have been defined", func() {
			c := goldigen.TypesConfiguration{
				Types: map[string]goldigen.TypeDefinition{},
			}
			Expect(c.Validate()).To(MatchError("no types have been defined: please define at least one type"))
		})

		It("should return an error if a type definition is missing a package", func() {
			c := goldigen.TypesConfiguration{
				Types: map[string]goldigen.TypeDefinition{
					"foo": goldigen.TypeDefinition{},
				},
			}
			Expect(c.Validate()).To(MatchError(`type definition of "foo" is missing the required "package" key`))
		})

		It("should return an error if a type definition is missing the factory method", func() {
			c := goldigen.TypesConfiguration{
				Types: map[string]goldigen.TypeDefinition{
					"foo": goldigen.TypeDefinition{
						Package: "foo/bar",
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type definition of "foo" is missing the required "factory" key`))
		})
	})

	Describe("retrieving all packages", func() {
		It("should return an empty list if no types were defined", func() {
			c := goldigen.TypesConfiguration{}
			Expect(c.Packages()).To(BeEmpty())
		})

		Context("with each type package appearing only once", func() {
			It("should return the packages alphabetically sorted", func() {
				c := goldigen.TypesConfiguration{
					Types: map[string]goldigen.TypeDefinition{
						"foo": goldigen.TypeDefinition{
							Package: "foo/test/package1",
						},
						"bar": goldigen.TypeDefinition{
							Package: "bar/test/package2",
						},
						"baz": goldigen.TypeDefinition{
							Package: "baz/test/package3",
						},
					},
				}

				Expect(c.Packages()).To(HaveLen(3))
				Expect(c.Packages()[0]).To(Equal("bar/test/package2"))
				Expect(c.Packages()[1]).To(Equal("baz/test/package3"))
				Expect(c.Packages()[2]).To(Equal("foo/test/package1"))
			})
		})

		Context("with packages appearing multiple times", func() {
			It("should return the packages", func() {
				c := goldigen.TypesConfiguration{
					Types: map[string]goldigen.TypeDefinition{
						"foo.1": goldigen.TypeDefinition{
							Package: "foo/test/package1",
						},
						"foo.2": goldigen.TypeDefinition{
							Package: "foo/test/package1",
						},
						"bar": goldigen.TypeDefinition{
							Package: "bar/test/package2",
						},
					},
				}
				Expect(c.Packages()).To(HaveLen(2))
				Expect(c.Packages()).To(ContainElement("foo/test/package1"))
				Expect(c.Packages()).To(ContainElement("bar/test/package2"))
			})
		})
	})
})
