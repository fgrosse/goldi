package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
	"fmt"
)

var _ = Describe("TypeConfigurator", func() {
	var container *goldi.Container

	BeforeEach(func() {
		container = goldi.NewContainer(goldi.NewTypeRegistry(), map[string]interface{}{})
	})

	Describe("Configure", func() {
		Context("when the type configurator has not been defined correctly", func() {
			It("should return an error", func() {
				someType := new(Foo)
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configurator.Configure(someType, container)).To(MatchError(`the configurator type "@configurator" has not been defined`))
			})
		})

		Context("when the type configurator is no struct or pointer to struct", func() {
			It("should return an error", func() {
				container.InjectInstance("configurator", 42)
				someType := new(Foo)
				configurator := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configurator.Configure(someType, container)).To(MatchError("the configurator instance is no struct or pointer to struct but a int"))
			})
		})

		Context("when the type configurator method does not exist", func() {
			It("should return an error", func() {
				someType := new(Foo)
				configurator := goldi.NewTypeConfigurator("configurator", "Fooobar")
				container.Register("configurator", goldi.NewInstanceType(&MyConfigurator{}))

				Expect(configurator.Configure(someType, container)).To(MatchError(`the configurator does not have a method "Fooobar"`))
			})
		})

		Context("when the type configurator has been defined properly", func() {
			var configurator *MyConfigurator
			BeforeEach(func() {
				configurator = &MyConfigurator{ConfiguredValue: "success!"}
				container.Register("configurator", goldi.NewInstanceType(configurator))
			})

			It("should return an error if the first argument is nil", func() {
				configuratorType := goldi.NewTypeConfigurator("configurator", "Configure")
				Expect(configuratorType.Configure(nil, container)).To(MatchError("can not configure nil"))
			})

			It("should call the requested function on the configurator", func() {
				someType := new(Foo)
				configuratorType := goldi.NewTypeConfigurator("configurator", "Configure")

				Expect(someType.Value).NotTo(Equal("success!"))
				Expect(configuratorType.Configure(someType, container)).To(Succeed())
				Expect(someType.Value).To(Equal("success!"))
			})

			It("should return an error if the configurator returned an error", func() {
				configurator.ReturnError = true

				someType := new(Foo)
				configuratorType := goldi.NewTypeConfigurator("configurator", "Configure")

				Expect(configuratorType.Configure(someType, container)).To(MatchError("this is the error message from the tests.MockTypeConfigurator"))
			})

			It("should return nil if the configurator returned nil", func() {
				configurator.ReturnError = false

				someType := new(Foo)
				configuratorType := goldi.NewTypeConfigurator("configurator", "Configure")

				Expect(configuratorType.Configure(someType, container)).To(BeNil())
			})
		})
	})
})

type MyConfigurator struct {
	ConfiguredValue string
	ReturnError bool
}

func (c *MyConfigurator) Configure(f *Foo) error {
	if c.ReturnError {
		return fmt.Errorf("this is the error message from the tests.MockTypeConfigurator")
	}

	f.Value = c.ConfiguredValue
	return nil
}
