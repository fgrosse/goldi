package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi/tests/testAPI"
)

var _ = Describe("TypeConfigurator", func() {
	var (
		config       = map[string]interface{}{}
		typeRegistry goldi.TypeRegistry
		resolver     *goldi.ParameterResolver
	)

	BeforeEach(func() {
		typeRegistry = goldi.NewTypeRegistry()
		resolver = goldi.NewParameterResolver(config, typeRegistry)
	})

	Describe("Configure", func() {
		Context("when the type configurator has not been defined correctly", func() {
			It("should return an error", func() {
				someType := testAPI.NewMockType()
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configurator.Configure(someType, resolver)).To(MatchError(`the referenced type "@configurator" has not been defined`))
			})
		})

		Context("when the type configurator is no struct or pointer to struct", func() {
			It("should return an error", func() {
				typeRegistry.InjectInstance("configurator", 42)
				someType := testAPI.NewMockType()
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configurator.Configure(someType, resolver)).To(MatchError("the configurator instance is no struct or pointer to struct but a int"))
			})
		})

		Context("when the type configurator method does not exist", func() {
			It("should return an error", func() {
				someType := testAPI.NewMockType()
				configurator := goldi.NewTypeConfigurator("configurator", "Fooobar")
				typeRegistry.RegisterType("configurator", testAPI.NewMockTypeConfigurator, "foobar")

				Expect(configurator.Configure(someType, resolver)).To(MatchError(`the configurator does not have a method "Fooobar"`))
			})
		})

		Context("when the type configurator has been defined properly", func() {
			BeforeEach(func() {
				typeRegistry.RegisterType("configurator", testAPI.NewMockTypeConfigurator, "foobar")
			})

			It("should return an error if the first argument is nil", func() {
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configurator.Configure(nil, resolver)).To(MatchError("can not configure nil"))
			})

			It("should call the requested function on the configurator", func() {
				someType := testAPI.NewMockType()
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")

				Expect(someType.StringParameter).NotTo(Equal("foobar"))
				Expect(configurator.Configure(someType, resolver)).To(Succeed())
				Expect(someType.StringParameter).To(Equal("foobar"))
			})

			It("should return an error if the configurator returned an error", func() {
				typeRegistry.InjectInstance("configurator", testAPI.NewFailingMockTypeConfigurator())

				someType := testAPI.NewMockType()
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")

				Expect(configurator.Configure(someType, resolver)).To(MatchError("this is the error message from the testAPI.MockTypeConfigurator"))
			})
		})
	})
})
