package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi/generator"
)

var _ = Describe("TypesConfiguration", func() {
	Describe("validation", func() {
		It("should return an error if no types have been defined", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{},
			}
			Expect(c.Validate()).To(MatchError("no types have been defined: please define at least one type"))
		})

		It("should return an error if a type definition is missing a package", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{},
				},
			}
			Expect(c.Validate()).To(MatchError(`type definition of "foo" is missing the required "package" key`))
		})

		It("should return an error if a type definition is missing the factory method", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{
						Package: "foo/bar",
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type definition of "foo" is missing the required "factory" key`))
		})

		It("should return an error if a type is an alias and contains a factory method", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{
						AliasForType: "bar",
						FactoryMethod: "NewFoo",
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type alias "foo" must not define a factory method`))
		})

		It("should return an error if a type is an alias and contains a package name", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{
						AliasForType: "bar",
						Package: "github.com/fgrosse/foo",
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type alias "foo" must not define a package name`))
		})

		It("should return an error if a type is an alias and contains a func", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{
						AliasForType: "bar",
						FuncName: "DoStuff",
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type alias "foo" must not define a func`))
		})

		It("should return an error if a type is an alias and contains arguments", func() {
			c := generator.TypesConfiguration{
				Types: map[string]generator.TypeDefinition{
					"foo": generator.TypeDefinition{
						AliasForType: "bar",
						RawArguments: []interface{}{"a", "b", "c"},
					},
				},
			}
			Expect(c.Validate()).To(MatchError(`type alias "foo" must not contain arguments`))
		})
	})

	Describe("retrieving all packages", func() {
		It("should return an empty list if no types were defined", func() {
			c := generator.TypesConfiguration{}
			Expect(c.Packages()).To(BeEmpty())
		})

		Context("with each type package appearing only once", func() {
			It("should return the packages alphabetically sorted", func() {
				c := generator.TypesConfiguration{
					Types: map[string]generator.TypeDefinition{
						"foo": generator.TypeDefinition{
							Package: "foo/test/package1",
						},
						"bar": generator.TypeDefinition{
							Package: "bar/test/package2",
						},
						"baz": generator.TypeDefinition{
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
				c := generator.TypesConfiguration{
					Types: map[string]generator.TypeDefinition{
						"foo.1": generator.TypeDefinition{
							Package: "foo/test/package1",
						},
						"foo.2": generator.TypeDefinition{
							Package: "foo/test/package1",
						},
						"bar": generator.TypeDefinition{
							Package: "bar/test/package2",
						},
					},
				}
				Expect(c.Packages()).To(HaveLen(2))
				Expect(c.Packages()).To(ContainElement("foo/test/package1"))
				Expect(c.Packages()).To(ContainElement("bar/test/package2"))
			})
		})

		Context("with packages from the arguments appearing in the configuration", func() {
			It("should return the packages", func() {
				c := generator.TypesConfiguration{
					Types: map[string]generator.TypeDefinition{
						"some_goldi_type": generator.TypeDefinition{
							Package: "github.com/fgrosse/goldi",
						},
					},
				}

				packages := c.Packages("github.com/fgrosse/goldi")
				Expect(packages).To(HaveLen(1))
				Expect(packages).To(ContainElement("github.com/fgrosse/goldi"))
			})
		})
	})

	PDescribe("short form of defining a type aliases", func() {
		// TODO
	})
})
