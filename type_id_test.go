package goldi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/goldi"
)

var _ = Describe("TypeID", func() {
	Describe("NewTypeID", func() {
		It("should panic if given an empty string", func() {
			Expect(func() {goldi.NewTypeID("")}).To(Panic())
		})
	})

	Describe("String", func() {
		It("should return the raw type if present", func() {
			t := goldi.TypeID{Raw: "@foo"}
			Expect(t.String()).To(Equal("@foo"))
		})

		It("should work if just the ID is set", func() {
			t := goldi.TypeID{ID: "foo"}
			Expect(t.String()).To(Equal("@foo"))
		})

		It("should use FuncReferenceMethod if it is not empty", func() {
			t := goldi.TypeID{ID: "foo", FuncReferenceMethod: "DoStuff"}
			Expect(t.String()).To(Equal("@foo::DoStuff"))
		})
	})
})
